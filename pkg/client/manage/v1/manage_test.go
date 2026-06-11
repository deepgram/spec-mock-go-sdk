// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Manage REST facade. Exercises the transport paths
// the generator wires up: multi-label path substitution, query params, a
// JSON request body on POST/PUT, response decoding, and the Token auth header —
// all against httptest, no live API.
package managev1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type capture struct {
	method string
	path   string
	query  string
	auth   string
	body   string
}

// server returns an httptest server that records the last request into *c and
// replies with the given JSON body.
func server(t *testing.T, c *capture, reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		c.method, c.path, c.query, c.auth, c.body = r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Get("Authorization"), string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, reply)
	}))
}

func TestListProjects(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"projects":[{"project_id":"p1","name":"Acme"}]}`)
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).ListProjects(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if c.method != "GET" || c.path != "/v1/projects" {
		t.Fatalf("got %s %s", c.method, c.path)
	}
	if c.auth != "Token k" {
		t.Fatalf("auth = %q", c.auth)
	}
	if len(out.Projects) != 1 || out.Projects[0].Name == nil || *out.Projects[0].Name != "Acme" {
		t.Fatalf("projects = %+v", out.Projects)
	}
}

func TestGetProjectKey_MultiLabelPath(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"item":{"api_key":{"api_key_id":"k1","comment":"ci"}}}`)
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).
		GetProjectKey(context.Background(), "proj1", "key1", nil)
	if err != nil {
		t.Fatalf("GetProjectKey: %v", err)
	}
	if c.path != "/v1/projects/proj1/keys/key1" {
		t.Fatalf("path = %q (multi-label substitution wrong)", c.path)
	}
	if out.Item == nil || out.Item.Api_key == nil || out.Item.Api_key.Api_key_id == nil || *out.Item.Api_key.Api_key_id != "k1" {
		t.Fatalf("item = %+v", out.Item)
	}
}

func TestCreateProjectKey_PostBody(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"api_key_id":"k2","key":"secret-once"}`)
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).CreateProjectKey(
		context.Background(), "proj1",
		map[string]any{"comment": "ci key", "scopes": []string{"member"}}, nil)
	if err != nil {
		t.Fatalf("CreateProjectKey: %v", err)
	}
	if c.method != "POST" || c.path != "/v1/projects/proj1/keys" {
		t.Fatalf("got %s %s", c.method, c.path)
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(c.body), &sent); err != nil {
		t.Fatalf("body not JSON: %q", c.body)
	}
	if sent["comment"] != "ci key" {
		t.Fatalf("body = %q", c.body)
	}
	if out.Key == nil || *out.Key != "secret-once" {
		t.Fatalf("key = %v", out.Key)
	}
}

func TestUpdateMemberScopes_Put(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"message":"ok"}`)
	defer srv.Close()

	_, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).UpdateMemberScopes(
		context.Background(), "proj1", "mem1",
		map[string]any{"scope": "member"}, nil)
	if err != nil {
		t.Fatalf("UpdateMemberScopes: %v", err)
	}
	if c.method != "PUT" || c.path != "/v1/projects/proj1/members/mem1/scopes" {
		t.Fatalf("got %s %s", c.method, c.path)
	}
}

func TestDeleteProject(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"message":"deleted"}`)
	defer srv.Close()

	_, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).DeleteProject(context.Background(), "proj1", nil)
	if err != nil {
		t.Fatalf("DeleteProject: %v", err)
	}
	if c.method != "DELETE" || c.path != "/v1/projects/proj1" {
		t.Fatalf("got %s %s", c.method, c.path)
	}
}

func TestGetUsageBreakdown_Query(t *testing.T) {
	var c capture
	srv := server(t, &c, `{"start":"a","end":"b","resolution":{"units":"hour","amount":1},"results":[]}`)
	defer srv.Close()

	_, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).GetUsageBreakdown(
		context.Background(), "proj1",
		&GetUsageBreakdownOptions{Grouping: "endpoint", Smart_format: true})
	if err != nil {
		t.Fatalf("GetUsageBreakdown: %v", err)
	}
	if c.path != "/v1/projects/proj1/usage/breakdown" {
		t.Fatalf("path = %q", c.path)
	}
	// both the grouping and the mixed-in filter must be on the wire
	if !contains(c.query, "grouping=endpoint") || !contains(c.query, "smart_format=true") {
		t.Fatalf("query = %q", c.query)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
