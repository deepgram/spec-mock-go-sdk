// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the multi-op REST facade: verifies a path-param GET,
// a query-param GET, and a POST-with-body op route correctly through the shared
// transport (path substitution, query binding, JSON body) and decode.
package agentsettingsv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentSettingsMultiOp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/v1/agent/settings/think/providers/anthropic":
			_, _ = w.Write([]byte(`{"id":"anthropic","name":"Anthropic"}`))
		case r.Method == "GET" && r.URL.Path == "/v1/agent/settings/think/providers":
			if r.URL.Query().Get("include") != "models" { // proves query binding reached the wire
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"providers":[{"id":"open_ai","name":"OpenAI"}]}`))
		case r.Method == "POST" && r.URL.Path == "/v1/agent/settings/validate":
			w.WriteHeader(http.StatusOK) // 200 empty = valid
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL), WithAPIKey("k"))
	ctx := context.Background()

	// Path-param GET — success implies {provider} was substituted into the URL.
	gp, err := c.GetThinkProvider(ctx, "anthropic", nil)
	if err != nil {
		t.Fatalf("GetThinkProvider: %v", err)
	}
	if gp.ID != "anthropic" {
		t.Fatalf("GetThinkProvider id = %v", gp.ID)
	}

	// Query-param GET — the server 400s unless include=models reached it.
	lp, err := c.ListThinkProviders(ctx, &ListThinkProvidersOptions{Include: "models"})
	if err != nil {
		t.Fatalf("ListThinkProviders: %v", err)
	}
	if len(lp.Providers) != 1 {
		t.Fatalf("ListThinkProviders providers = %d", len(lp.Providers))
	}

	// POST with JSON body.
	if _, err := c.ValidateSettings(ctx, map[string]any{"type": "Settings"}, nil); err != nil {
		t.Fatalf("ValidateSettings: %v", err)
	}
}
