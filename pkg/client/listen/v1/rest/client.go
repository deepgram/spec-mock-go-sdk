// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

const DefaultBaseURL = "https://api.deepgram.com"

// Client is the Listen REST API client. It holds credentials and
// transport configuration shared across all calls to FromURL,
// FromFile, and FromStream.
//
// The zero value of Client is not usable. Construct via New or
// NewWithDefaults.
type Client struct {
	apiKey      string
	accessToken string
	baseURL     string
	httpClient  *http.Client
}

// New constructs a Client with explicit credentials. Either apiKey
// or accessToken must be non-empty for requests to authenticate;
// accessToken takes precedence when both are set.
func New(apiKey, accessToken string) *Client {
	return &Client{
		apiKey:      apiKey,
		accessToken: accessToken,
		baseURL:     DefaultBaseURL,
		httpClient:  http.DefaultClient,
	}
}

// NewWithDefaults reads DEEPGRAM_ACCESS_TOKEN and DEEPGRAM_API_KEY
// from the environment.
func NewWithDefaults() *Client {
	return New(
		os.Getenv("DEEPGRAM_API_KEY"),
		os.Getenv("DEEPGRAM_ACCESS_TOKEN"),
	)
}

// WithBaseURL returns a copy of the Client pointed at the given base URL.
func (c *Client) WithBaseURL(url string) *Client {
	out := *c
	out.baseURL = url
	return &out
}

// WithHTTPClient returns a copy of the Client using the given http.Client.
func (c *Client) WithHTTPClient(hc *http.Client) *Client {
	out := *c
	out.httpClient = hc
	return &out
}

// FromURL transcribes audio at a remote URL. The URL is delivered to
// the Deepgram API as a JSON body; the API fetches and transcribes.
func (c *Client) FromURL(ctx context.Context, audioURL string, opts *PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	body := strings.NewReader(`{"url":` + jsonQuote(audioURL) + `}`)
	return c.invoke(ctx, opts, "application/json", body)
}

// FromFile transcribes audio at a local file path. The file is
// streamed as the request body with the supplied contentType (or
// "audio/*" when empty).
func (c *Client) FromFile(ctx context.Context, path, contentType string, opts *PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return c.FromStream(ctx, f, contentType, opts)
}

// FromStream transcribes audio from an arbitrary io.Reader.
// contentType defaults to "audio/*" if empty.
func (c *Client) FromStream(ctx context.Context, r io.Reader, contentType string, opts *PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	if contentType == "" {
		contentType = "audio/*"
	}
	return c.invoke(ctx, opts, contentType, r)
}

func (c *Client) invoke(ctx context.Context, opts *PreRecordedTranscriptionOptions, contentType string, body io.Reader) (*PreRecordedResponse, error) {
	input := optionsToTranscribeInput(opts)
	out, err := httptransport.Invoke[spectypes.TranscribeInput, spectypes.TranscribeOutput](
		ctx,
		c.httpClient,
		c.baseURL,
		spectypes.TranscribeRoute,
		&listenAuthenticator{
			apiKey:      c.apiKey,
			accessToken: c.accessToken,
			contentType: contentType,
		},
		input,
		body,
	)
	if err != nil {
		return nil, err
	}
	return convertTranscribeOutput(out), nil
}

// listenAuthenticator applies Authorization and Content-Type headers
// to the outgoing request. It piggybacks Content-Type onto the
// Authenticator hook because the generic httptransport.Invoke does
// not expose a separate header-setting hook for the body payload.
type listenAuthenticator struct {
	apiKey      string
	accessToken string
	contentType string
}

func (a *listenAuthenticator) Apply(r *http.Request) error {
	if a.contentType != "" {
		r.Header.Set("Content-Type", a.contentType)
	}
	switch {
	case a.accessToken != "":
		r.Header.Set("Authorization", "Bearer "+a.accessToken)
	case a.apiKey != "":
		r.Header.Set("Authorization", "Token "+a.apiKey)
	default:
		return errors.New("listen rest: no credentials; set DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN, or pass them to New(...)")
	}
	return nil
}

// jsonQuote produces a JSON string literal for s. Used to construct
// the {"url": "..."} body for FromURL without pulling in encoding/json.
func jsonQuote(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
