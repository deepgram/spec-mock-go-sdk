// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Read text-intelligence facade: a single-op
// content-type-union product. FromText sends text/plain, FromURL sends a
// JSON {"url": ...} body; both POST to /v1/read and wire the typed options as
// query params. Verifies the request shape and that the JSON response decodes
// into the idiomatic ReadResponse.
package readv1

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFromText_WiresPlainTextAndQuery(t *testing.T) {
	var gotMethod, gotPath, gotCT, gotAuth, gotBody string
	var gotQuery map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotMethod, gotPath = r.Method, r.URL.Path
		gotCT, gotAuth, gotBody = r.Header.Get("Content-Type"), r.Header.Get("Authorization"), string(b)
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"metadata":{"request_id":"req-text"}}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromText(
		context.Background(), "analyze me", &ReadOptions{Language: "en", Sentiment: true})
	if err != nil {
		t.Fatalf("FromText: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/v1/read" {
		t.Fatalf("got %s %s, want POST /v1/read", gotMethod, gotPath)
	}
	if gotCT != "text/plain" {
		t.Fatalf("Content-Type = %q, want text/plain", gotCT)
	}
	if gotAuth != "Token k" {
		t.Fatalf("auth = %q", gotAuth)
	}
	if gotBody != "analyze me" {
		t.Fatalf("body = %q, want %q", gotBody, "analyze me")
	}
	if q := gotQuery["language"]; len(q) != 1 || q[0] != "en" {
		t.Fatalf("language query = %v, want [en]", q)
	}
	if q := gotQuery["sentiment"]; len(q) != 1 || q[0] != "true" {
		t.Fatalf("sentiment query = %v, want [true]", q)
	}
	if out.Metadata == nil || out.Metadata.RequestID != "req-text" {
		t.Fatalf("metadata.request_id = %v, want req-text", out.Metadata)
	}
}

func TestFromURL_WiresJSONBody(t *testing.T) {
	var gotCT, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotCT, gotBody = r.Header.Get("Content-Type"), string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"metadata":{"request_id":"req-url"}}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromURL(
		context.Background(), "https://example.invalid/article.txt", nil)
	if err != nil {
		t.Fatalf("FromURL: %v", err)
	}
	if gotCT != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", gotCT)
	}
	if gotBody != `{"url":"https://example.invalid/article.txt"}` {
		t.Fatalf("body = %q", gotBody)
	}
	if out.Metadata == nil || out.Metadata.RequestID != "req-url" {
		t.Fatalf("metadata.request_id = %v, want req-url", out.Metadata)
	}
}
