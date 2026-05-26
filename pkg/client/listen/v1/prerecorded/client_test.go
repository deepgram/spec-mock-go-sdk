// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package prerecordedv1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func TestFromURL_AdditionalQueryParamsOverrideTypedField(t *testing.T) {
	var capturedQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"req-123"}`))
	}))
	defer server.Close()

	client := prerecorded.New(
		prerecorded.WithCredentials("test-api-key", ""),
		prerecorded.WithBaseURL(server.URL),
	)
	_, err := client.FromURL(context.Background(), "https://example.invalid/audio.wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model:                 "nova-3",
		AdditionalQueryParams: url.Values{"model": []string{"nova-2-meeting"}},
	})
	if err != nil {
		t.Fatalf("FromURL returned error: %v", err)
	}

	if got := capturedQuery["model"]; len(got) != 1 || got[0] != "nova-2-meeting" {
		t.Fatalf("model query = %v, want [nova-2-meeting]", got)
	}
}

func TestNew_NoOptions_UsesDefaults(t *testing.T) {
	t.Setenv("DEEPGRAM_API_KEY", "env-api-key")
	t.Setenv("DEEPGRAM_ACCESS_TOKEN", "")

	client := prerecorded.New()
	if got := clientStringField(t, client, "apiKey"); got != "env-api-key" {
		t.Fatalf("apiKey = %q, want env-api-key", got)
	}
	if got := clientStringField(t, client, "baseURL"); got != prerecorded.DefaultBaseURL {
		t.Fatalf("baseURL = %q, want %q", got, prerecorded.DefaultBaseURL)
	}
	if clientPtrField(t, client, "httpClient") == 0 {
		t.Fatal("httpClient was nil")
	}
}

func TestNew_WithCredentials(t *testing.T) {
	client := prerecorded.New(prerecorded.WithCredentials("api-key", "access-token"))
	if got := clientStringField(t, client, "apiKey"); got != "api-key" {
		t.Fatalf("apiKey = %q, want api-key", got)
	}
	if got := clientStringField(t, client, "accessToken"); got != "access-token" {
		t.Fatalf("accessToken = %q, want access-token", got)
	}
}

func TestNew_WithBaseURL(t *testing.T) {
	client := prerecorded.New(prerecorded.WithBaseURL("https://api.example.invalid"))
	if got := clientStringField(t, client, "baseURL"); got != "https://api.example.invalid" {
		t.Fatalf("baseURL = %q, want custom URL", got)
	}
}

func TestNew_WithHTTPClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	client := prerecorded.New(prerecorded.WithHTTPClient(httpClient))
	if got, want := clientPtrField(t, client, "httpClient"), reflect.ValueOf(httpClient).Pointer(); got != want {
		t.Fatalf("httpClient pointer = %#x, want %#x", got, want)
	}
}

func TestNew_WithSageMakerTransport(t *testing.T) {
	awsClient := sagemakerruntime.New(sagemakerruntime.Options{Region: "us-east-1"})
	client := prerecorded.New(prerecorded.WithSageMakerTransport(awsClient, "endpoint-name"))
	if clientInterfaceFieldIsNil(t, client, "transport") {
		t.Fatal("transport was nil")
	}
}

func clientStringField(t *testing.T, client *prerecorded.Client, name string) string {
	t.Helper()
	return reflect.ValueOf(client).Elem().FieldByName(name).String()
}

func clientPtrField(t *testing.T, client *prerecorded.Client, name string) uintptr {
	t.Helper()
	return reflect.ValueOf(client).Elem().FieldByName(name).Pointer()
}

func clientInterfaceFieldIsNil(t *testing.T, client *prerecorded.Client, name string) bool {
	t.Helper()
	return reflect.ValueOf(client).Elem().FieldByName(name).IsNil()
}
