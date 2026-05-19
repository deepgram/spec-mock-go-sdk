// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func TestFromURL_HappyPath(t *testing.T) {
	var capturedAuth, capturedContentType string
	var capturedQuery url.Values
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedContentType = r.Header.Get("Content-Type")
		capturedQuery = r.URL.Query()
		capturedBody, _ = io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"request_id": "abc-123",
			"metadata": {
				"request_id": "abc-123",
				"transaction_key": "deprecated",
				"sha256": "xyz",
				"created": "2026-05-19T12:00:00Z",
				"duration": 4.2,
				"channels": 1,
				"model_info": {}
			},
			"results": {
				"channels": [{
					"alternatives": [{
						"transcript": "hello world",
						"confidence": 0.98,
						"words": []
					}]
				}]
			}
		}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)

	resp, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&PreRecordedTranscriptionOptions{
			Model:       "nova-3",
			Language:    "en-US",
			Punctuate:   true,
			SmartFormat: true,
		})
	if err != nil {
		t.Fatalf("FromURL returned error: %v", err)
	}

	if capturedAuth != "Token test-api-key" {
		t.Errorf("auth header: got %q, want %q", capturedAuth, "Token test-api-key")
	}
	if capturedContentType != "application/json" {
		t.Errorf("content-type: got %q, want %q", capturedContentType, "application/json")
	}
	if got := capturedQuery.Get("model"); got != "nova-3" {
		t.Errorf("?model: got %q, want %q", got, "nova-3")
	}
	if got := capturedQuery.Get("language"); got != "en-US" {
		t.Errorf("?language: got %q, want %q", got, "en-US")
	}
	if got := capturedQuery.Get("punctuate"); got != "true" {
		t.Errorf("?punctuate: got %q, want %q", got, "true")
	}
	if got := capturedQuery.Get("smart_format"); got != "true" {
		t.Errorf("?smart_format: got %q, want %q", got, "true")
	}

	var bodyJSON map[string]any
	if err := json.Unmarshal(capturedBody, &bodyJSON); err != nil {
		t.Fatalf("body was not JSON: %v\nraw: %s", err, capturedBody)
	}
	if bodyJSON["url"] != "https://dpgr.am/spacewalk.wav" {
		t.Errorf(`body.url: got %v, want %q`, bodyJSON["url"], "https://dpgr.am/spacewalk.wav")
	}

	if resp.RequestID != "abc-123" {
		t.Errorf("resp.RequestID: got %q, want %q", resp.RequestID, "abc-123")
	}
	transcript := resp.Results.Channels[0].Alternatives[0].Transcript
	if transcript != "hello world" {
		t.Errorf("transcript: got %q, want %q", transcript, "hello world")
	}
}

func TestFromURL_AccessTokenPrecedence(t *testing.T) {
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"x","results":{"channels":[{"alternatives":[{"transcript":""}]}]}}`))
	}))
	defer server.Close()

	client := New("the-api-key", "the-access-token").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(), "https://example.invalid/audio.wav", nil)
	if err != nil {
		t.Fatalf("FromURL returned error: %v", err)
	}
	if capturedAuth != "Bearer the-access-token" {
		t.Errorf("auth header: got %q, want %q (AccessToken should win over APIKey)",
			capturedAuth, "Bearer the-access-token")
	}
}

func TestFromURL_NoCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("server should not have been hit; auth should fail before send")
	}))
	defer server.Close()

	client := New("", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(), "https://example.invalid/audio.wav", nil)
	if err == nil {
		t.Fatal("FromURL with no credentials should return an error")
	}
	if !strings.Contains(err.Error(), "no credentials") {
		t.Errorf("error message: got %q, want substring %q", err.Error(), "no credentials")
	}
}

func TestFromURL_HTTPErrorIsTypedAndDiscriminable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Dg-Request-Id", "req-deadbeef")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err_code":"Bad Request","err_msg":"invalid model"}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&PreRecordedTranscriptionOptions{Model: "bogus"})
	if err == nil {
		t.Fatal("FromURL with 400 response should return an error")
	}

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err should be assertable to *HTTPError via errors.As; got %T: %v", err, err)
	}

	if httpErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode: got %d, want %d", httpErr.StatusCode, http.StatusBadRequest)
	}
	if !strings.Contains(string(httpErr.Body), "invalid model") {
		t.Errorf("Body: got %q, want substring %q", string(httpErr.Body), "invalid model")
	}
	if got := httpErr.Headers.Get("X-Dg-Request-Id"); got != "req-deadbeef" {
		t.Errorf("Headers[X-Dg-Request-Id]: got %q, want %q", got, "req-deadbeef")
	}
	if httpErr.Method != "POST" {
		t.Errorf("Method: got %q, want %q", httpErr.Method, "POST")
	}
}

func TestFromURL_HTTPErrorTypedPopulatedFor400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("dg-error", "invalid model: bogus")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err_code":"INVALID_QUERY_PARAMETER","err_msg":"invalid model: bogus","request_id":"01HXYZ-deadbeef"}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&PreRecordedTranscriptionOptions{Model: "bogus"})
	if err == nil {
		t.Fatal("FromURL with 400 response should return an error")
	}

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err should be assertable to *HTTPError; got %T: %v", err, err)
	}
	if httpErr.Typed == nil {
		t.Fatal("HTTPError.Typed should be populated for a declared 400 error")
	}

	var queryErr *spectypes.InvalidQueryParameterError
	if !errors.As(httpErr.Typed, &queryErr) {
		t.Fatalf("Typed should be *InvalidQueryParameterError; got %T", httpErr.Typed)
	}

	if queryErr.ErrCode == nil || *queryErr.ErrCode != "INVALID_QUERY_PARAMETER" {
		t.Errorf("ErrCode: got %v, want \"INVALID_QUERY_PARAMETER\"", queryErr.ErrCode)
	}
	if queryErr.ErrMsg == nil || *queryErr.ErrMsg != "invalid model: bogus" {
		t.Errorf("ErrMsg: got %v, want \"invalid model: bogus\"", queryErr.ErrMsg)
	}
	if queryErr.RequestId == nil || *queryErr.RequestId != "01HXYZ-deadbeef" {
		t.Errorf("RequestId: got %v, want \"01HXYZ-deadbeef\"", queryErr.RequestId)
	}
	if queryErr.DgError == nil || *queryErr.DgError != "invalid model: bogus" {
		t.Errorf("DgError (from header): got %v, want \"invalid model: bogus\"", queryErr.DgError)
	}
}

func TestFromURL_HTTPErrorTypedPopulatedFor429WithRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "60")
		w.Header().Set("dg-error", "rate limit hit")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"err_code":"RATE_LIMITED","err_msg":"rate limit hit","request_id":"abc"}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav", nil)
	if err == nil {
		t.Fatal("FromURL with 429 response should return an error")
	}

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err should be assertable to *HTTPError; got %T", err)
	}

	var rateLimited *spectypes.RateLimitedError
	if !errors.As(httpErr.Typed, &rateLimited) {
		t.Fatalf("Typed should be *RateLimitedError; got %T", httpErr.Typed)
	}
	if rateLimited.RetryAfter == nil || *rateLimited.RetryAfter != "60" {
		t.Errorf("RetryAfter: got %v, want \"60\"", rateLimited.RetryAfter)
	}
}

func TestFromURL_ExtraQueryParamsAdded(t *testing.T) {
	var capturedQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"x","results":{"channels":[{"alternatives":[{"transcript":""}]}]}}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://example.invalid/audio.wav",
		&PreRecordedTranscriptionOptions{
			Model: "nova-3",
			Extra: url.Values{
				"experimental_feature":   []string{"true"},
				"custom_tag":             []string{"a", "b"},
			},
		})
	if err != nil {
		t.Fatalf("FromURL returned error: %v", err)
	}

	if got := capturedQuery.Get("model"); got != "nova-3" {
		t.Errorf("typed field ?model: got %q, want %q", got, "nova-3")
	}
	if got := capturedQuery.Get("experimental_feature"); got != "true" {
		t.Errorf("?experimental_feature: got %q, want %q", got, "true")
	}
	vs := capturedQuery["custom_tag"]
	if len(vs) != 2 || vs[0] != "a" || vs[1] != "b" {
		t.Errorf("?custom_tag: got %v, want [a b]", vs)
	}
}

func TestFromURL_ExtraOverridesTypedField(t *testing.T) {
	var capturedQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"x","results":{"channels":[{"alternatives":[{"transcript":""}]}]}}`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://example.invalid/audio.wav",
		&PreRecordedTranscriptionOptions{
			Model: "nova-3",
			Extra: url.Values{"model": []string{"nova-2-meeting"}},
		})
	if err != nil {
		t.Fatalf("FromURL returned error: %v", err)
	}

	vs := capturedQuery["model"]
	if len(vs) != 1 || vs[0] != "nova-2-meeting" {
		t.Errorf("?model: got %v, want [nova-2-meeting] (Extra should win)", vs)
	}
}

func TestFromURL_HTTPErrorTypedNilForUndecodableBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`<html><body>502 Bad Gateway</body></html>`))
	}))
	defer server.Close()

	client := New("test-api-key", "").WithBaseURL(server.URL)
	_, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav", nil)
	if err == nil {
		t.Fatal("FromURL with 502 response should return an error")
	}

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err should be assertable to *HTTPError; got %T", err)
	}
	if httpErr.StatusCode != http.StatusBadGateway {
		t.Errorf("StatusCode: got %d, want %d", httpErr.StatusCode, http.StatusBadGateway)
	}
	if httpErr.Typed != nil {
		t.Errorf("Typed should be nil for unrecognised status (502 not declared); got %T", httpErr.Typed)
	}
	if !strings.Contains(string(httpErr.Body), "502 Bad Gateway") {
		t.Errorf("Body should still contain raw response; got %q", string(httpErr.Body))
	}
}
