// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the prerecorded client implementation for the Deepgram API.
//
// As of Phase 4 of the endpoint-agnostic transport migration, DoFile, DoStream,
// and DoURL all funnel through the generated transport/http.Invoke primitive
// with the TranscribeRoute metadata emitted in api/types/listen_route.go.
// Query-param assembly, header binding, auth, and response decoding live in
// the generic primitive; this client owns the spec-specific bits — body
// preparation (raw audio bytes vs JSON URL envelope) and Content-Type
// negotiation between the two payload variants of TranscribeRequestBody.
package restv1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	klog "k8s.io/klog/v2"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	common "github.com/deepgram/spec-mock-go-sdk/pkg/client/common/v1"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces/v1"
)

type urlSource struct {
	URL string `json:"url"`
}

// NewWithDefaults creates a new prerecorded client with all default options.
// The Deepgram API key is read from DEEPGRAM_API_KEY in the environment.
func NewWithDefaults() *Client {
	return New("", &interfaces.ClientOptions{})
}

// New creates a new prerecorded client. apiKey overrides any value already
// in options.APIKey; options carries host, version, path, and other knobs.
func New(apiKey string, options *interfaces.ClientOptions) *Client {
	if apiKey != "" {
		options.APIKey = apiKey
	}
	if err := options.Parse(); err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil
	}
	return &Client{common.NewREST(apiKey, options)}
}

// DoFile transcribes a local audio file. It opens the file and streams the
// bytes through DoStream.
func (c *Client) DoFile(ctx context.Context, filePath string, req *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	if len(req.Keyterm) > 0 && !strings.HasPrefix(req.Model, "nova-3") {
		klog.V(1).Info("Keyterms are only supported with nova-3 models.")
		return nil
	}

	info, err := os.Stat(filePath)
	if err != nil || errors.Is(err, os.ErrNotExist) {
		klog.V(1).Infof("File %s does not exist. Err: %v\n", filePath, err)
		return err
	}
	if info.IsDir() || info.Size() == 0 {
		klog.V(1).Infof("%s is empty or a directory\n", filePath)
		return ErrInvalidInput
	}

	file, err := os.Open(filePath)
	if err != nil {
		klog.V(1).Infof("os.Open(%s) failed. Err: %v\n", filePath, err)
		return err
	}
	defer file.Close()
	return c.DoStream(ctx, file, req, resBody)
}

// DoStream transcribes an audio byte stream. The stream is sent as the request
// body with Content-Type: application/octet-stream; query/header binding is
// handled by httptransport.Invoke from the TranscribeRoute metadata.
func (c *Client) DoStream(ctx context.Context, src io.Reader, options *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	if len(options.Keyterm) > 0 && !strings.HasPrefix(options.Model, "nova-3") {
		klog.V(1).Info("Keyterms are only supported with nova-3 models.")
		return nil
	}
	contentType := "application/octet-stream"
	return c.invoke(ctx, options, contentType, src, resBody)
}

// DoURL transcribes an audio file fetched from a remote URL. The URL is JSON-
// encoded into a {"url": "..."} envelope and sent with
// Content-Type: application/json.
func (c *Client) DoURL(ctx context.Context, audioURL string, options *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	if !IsURL(audioURL) {
		klog.V(1).Infof("Invalid URL: %s\n", audioURL)
		return ErrInvalidInput
	}
	body, err := json.Marshal(urlSource{URL: audioURL})
	if err != nil {
		klog.V(1).Infof("json.Marshal(urlSource) failed. Err: %v\n", err)
		return err
	}
	return c.invoke(ctx, options, "application/json", bytes.NewReader(body), resBody)
}

// invoke is the shared inner path that converts options to spectypes.TranscribeInput,
// fires httptransport.Invoke with TranscribeRoute metadata, and copies the typed
// *TranscribeOutput into resBody (the legacy interface{} sink kept for API
// compatibility with sendAudio in prerecorded.go).
func (c *Client) invoke(
	ctx context.Context,
	options *interfaces.PreRecordedTranscriptionOptions,
	contentType string,
	body io.Reader,
	resBody interface{},
) error {
	input := optionsToTranscribeInput(options)
	input.ContentType = &contentType

	out, err := httptransport.Invoke[spectypes.TranscribeInput, spectypes.TranscribeOutput](
		ctx, &c.HTTPClient.Client,
		c.baseURL(),
		spectypes.TranscribeRoute,
		c.authenticator(),
		input, body,
	)
	if err != nil {
		klog.V(1).Infof("httptransport.Invoke failed. Err: %v\n", err)
		return err
	}
	return decodeIntoResp(out, resBody)
}

// baseURL returns the scheme + host the Invoke primitive prepends to route.Path.
// If options.Host already carries a scheme prefix, it's returned verbatim;
// otherwise https:// is assumed.
func (c *Client) baseURL() string {
	host := c.Options.Host
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return host
	}
	return "https://" + host
}

// authenticator returns an httptransport.Authenticator that writes the
// Authorization header from options' AccessToken (Bearer) or APIKey (token).
func (c *Client) authenticator() httptransport.Authenticator {
	token, isBearer := c.Options.GetAuthToken()
	if token == "" {
		return nil
	}
	return &tokenAuth{token: token, bearer: isBearer}
}

type tokenAuth struct {
	token  string
	bearer bool
}

func (a *tokenAuth) Apply(req *http.Request) error {
	if a.bearer {
		req.Header.Set("Authorization", "Bearer "+a.token)
	} else {
		req.Header.Set("Authorization", "token "+a.token)
	}
	return nil
}

// decodeIntoResp copies src into dst. Fast path for *spectypes.TranscribeOutput;
// JSON-roundtrip fallback for other concrete types (e.g. tests passing custom
// shape). Preserves the resBody interface{} contract of the legacy Do* methods.
func decodeIntoResp(src *spectypes.TranscribeOutput, dst interface{}) error {
	if d, ok := dst.(*spectypes.TranscribeOutput); ok {
		*d = *src
		return nil
	}
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// IsURL returns true if s parses as a URL with a scheme and host.
func IsURL(s string) bool {
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return false
	}
	return len(s) > len("https://")
}
