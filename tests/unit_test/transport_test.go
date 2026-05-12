// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package deepgramtest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	smithyendpoints "github.com/aws/smithy-go/endpoints"

	sm "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	wrtc "github.com/deepgram/spec-mock-go-sdk/api/transport/webrtc"
)

type fixedEndpointResolver struct{ URL string }

func (r *fixedEndpointResolver) ResolveEndpoint(_ context.Context, _ sagemakerruntime.EndpointParameters) (smithyendpoints.Endpoint, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{URI: *u}, nil
}

func newSageMakerClientPointingAt(server *httptest.Server) *sagemakerruntime.Client {
	return sagemakerruntime.New(sagemakerruntime.Options{
		Region:             "us-east-1",
		Credentials:        aws.AnonymousCredentials{},
		EndpointResolverV2: &fixedEndpointResolver{URL: server.URL},
		HTTPClient:         server.Client(),
	})
}

func Test_SM_InvokeTranscribe_RoundtripsBodies(t *testing.T) {
	var received []byte
	var receivedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		received, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "sm-abc-123",
			"metadata": {
				"transaction_key": "deprecated",
				"request_id": "sm-abc-123",
				"sha256": "fakehash",
				"created": "2026-05-12T12:00:00Z",
				"duration": 1.0,
				"channels": 1,
				"model_info": {}
			},
			"results": {
				"channels": [{
					"alternatives": [{"transcript": "from sagemaker", "confidence": 0.9, "words": [], "languages": []}]
				}]
			}
		}`))
	}))
	defer server.Close()

	client := newSageMakerClientPointingAt(server)
	model := "nova-3"
	input := &spectypes.TranscribeInput{Model: &model}

	out, err := sm.InvokeTranscribe(context.Background(), client, "my-listen-endpoint", input)
	if err != nil {
		t.Fatalf("InvokeTranscribe: %v", err)
	}
	if out == nil || out.RequestId == nil || *out.RequestId != "sm-abc-123" {
		t.Fatalf("RequestId = %v, want sm-abc-123", out.RequestId)
	}
	if out.Results == nil || len(out.Results.Channels) == 0 {
		t.Fatal("Results.Channels empty")
	}
	transcript := out.Results.Channels[0].Alternatives[0].Transcript
	if transcript == nil || *transcript != "from sagemaker" {
		t.Fatalf("Transcript = %v, want 'from sagemaker'", transcript)
	}
	if receivedContentType != "application/json" {
		t.Fatalf("server received Content-Type = %q, want application/json", receivedContentType)
	}
	var roundtripped map[string]any
	if err := json.Unmarshal(received, &roundtripped); err != nil {
		t.Fatalf("server received non-JSON body: %v", err)
	}
}

func Test_SM_BidiStream_ReturnsPendingSDKError(t *testing.T) {
	_, err := sm.OpenListenLiveStream(context.Background(), "endpoint-name")
	if err == nil {
		t.Fatal("expected pending-SDK error, got nil")
	}
	if !strings.Contains(err.Error(), "pending aws-sdk-go-v2") {
		t.Fatalf("error message = %q, want it to mention pending aws-sdk-go-v2 support", err.Error())
	}
}

func Test_SM_BidiStream_StubMethods_Error(t *testing.T) {
	s := &sm.ListenLiveStream{}
	if err := s.Send(&spectypes.ClientStreamMemberCloseStream{}); err == nil {
		t.Fatal("Send: expected error, got nil")
	}
	if _, err := s.Recv(); err == nil {
		t.Fatal("Recv: expected error, got nil")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: unexpected error: %v", err)
	}
}

func Test_WebRTC_Stream_ReturnsNotImplementedError(t *testing.T) {
	_, err := wrtc.OpenListenLiveStream(context.Background())
	if err == nil {
		t.Fatal("expected not-implemented error, got nil")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("error message = %q, want it to mention 'not implemented'", err.Error())
	}
}

func Test_WebRTC_Stream_StubMethods_Error(t *testing.T) {
	s := &wrtc.ListenLiveStream{}
	if err := s.Send(&spectypes.ClientStreamMemberCloseStream{}); err == nil {
		t.Fatal("Send: expected error, got nil")
	}
	if _, err := s.Recv(); err == nil {
		t.Fatal("Recv: expected error, got nil")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: unexpected error: %v", err)
	}
}
