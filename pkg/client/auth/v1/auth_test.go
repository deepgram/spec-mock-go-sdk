// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Auth token-grant facade: a single-op REST service
// that is NOT a content-type-union product, so it routes through the generic
// multi-op emitter (POST + JSON body) rather than FromText/FromURL.
package authv1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGrantToken(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotMethod, gotPath, gotAuth, gotBody = r.Method, r.URL.Path, r.Header.Get("Authorization"), string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"tok-123","expires_in":30}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).GrantToken(
		context.Background(), map[string]any{"ttl_seconds": 30}, nil)
	if err != nil {
		t.Fatalf("GrantToken: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/v1/auth/grant" {
		t.Fatalf("got %s %s", gotMethod, gotPath)
	}
	if gotAuth != "Token k" {
		t.Fatalf("auth = %q", gotAuth)
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(gotBody), &sent); err != nil {
		t.Fatalf("body not JSON: %q", gotBody)
	}
	if out.Access_token == nil || *out.Access_token != "tok-123" {
		t.Fatalf("access_token = %v", out.Access_token)
	}
	if out.Expires_in == nil || *out.Expires_in != 30 {
		t.Fatalf("expires_in = %v", out.Expires_in)
	}
}
