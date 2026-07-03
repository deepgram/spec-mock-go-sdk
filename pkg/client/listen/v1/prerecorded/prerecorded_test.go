// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Listen prerecorded (batch) facade: a single-op
// content-type-union product whose raw variant is audio bytes. FromURL sends a
// JSON {"url": ...} body; FromFile/FromStream send the raw audio with the
// caller's content-type. All POST to /v1/listen and wire typed options as
// query params. Verifies the request shape, the audio-body path, and that
// AdditionalQueryParams override a typed field.
package prerecordedv1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestFromURL_WiresJSONBodyAndQuery(t *testing.T) {
	var gotMethod, gotPath, gotCT, gotAuth, gotBody string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotMethod, gotPath = r.Method, r.URL.Path
		gotCT, gotAuth, gotBody = r.Header.Get("Content-Type"), r.Header.Get("Authorization"), string(b)
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"request_id":"req-url"}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromURL(
		context.Background(), "https://example.invalid/audio.wav",
		&PreRecordedTranscriptionOptions{Encoding: "linear16", Diarize: true})
	if err != nil {
		t.Fatalf("FromURL: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/v1/listen" {
		t.Fatalf("got %s %s, want POST /v1/listen", gotMethod, gotPath)
	}
	if gotCT != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", gotCT)
	}
	if gotAuth != "Token k" {
		t.Fatalf("auth = %q", gotAuth)
	}
	if gotBody != `{"url":"https://example.invalid/audio.wav"}` {
		t.Fatalf("body = %q", gotBody)
	}
	if q := gotQuery["encoding"]; len(q) != 1 || q[0] != "linear16" {
		t.Fatalf("encoding query = %v, want [linear16]", q)
	}
	if q := gotQuery["diarize"]; len(q) != 1 || q[0] != "true" {
		t.Fatalf("diarize query = %v, want [true]", q)
	}
	if out.RequestID != "req-url" {
		t.Fatalf("request_id = %q, want req-url", out.RequestID)
	}
}

func TestFromStream_WiresAudioBody(t *testing.T) {
	audio := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x01, 0x02, 0x03} // fake WAV bytes
	var gotCT string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		gotCT = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"request_id":"req-stream"}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromStream(
		context.Background(), bytes.NewReader(audio), "audio/wav", nil)
	if err != nil {
		t.Fatalf("FromStream: %v", err)
	}
	if gotCT != "audio/wav" {
		t.Fatalf("Content-Type = %q, want audio/wav", gotCT)
	}
	if !bytes.Equal(gotBody, audio) {
		t.Fatalf("body = %v, want raw audio %v", gotBody, audio)
	}
	if out.RequestID != "req-stream" {
		t.Fatalf("request_id = %q, want req-stream", out.RequestID)
	}
}

func TestFromStream_DefaultsContentType(t *testing.T) {
	var gotCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCT = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()

	_, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromStream(
		context.Background(), bytes.NewReader([]byte{0x00}), "", nil)
	if err != nil {
		t.Fatalf("FromStream: %v", err)
	}
	if gotCT != "audio/*" {
		t.Fatalf("empty content-type defaulted to %q, want audio/*", gotCT)
	}
}

func TestFromFile_StreamsFileContents(t *testing.T) {
	audio := []byte{0x4f, 0x67, 0x67, 0x53, 0xde, 0xad, 0xbe, 0xef} // fake OGG bytes
	dir := t.TempDir()
	path := filepath.Join(dir, "clip.ogg")
	if err := os.WriteFile(path, audio, 0o600); err != nil {
		t.Fatalf("write temp audio: %v", err)
	}

	var gotCT string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		gotCT = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"request_id":"req-file"}`)
	}))
	defer srv.Close()

	out, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromFile(
		context.Background(), path, "audio/ogg", nil)
	if err != nil {
		t.Fatalf("FromFile: %v", err)
	}
	if gotCT != "audio/ogg" {
		t.Fatalf("Content-Type = %q, want audio/ogg", gotCT)
	}
	if !bytes.Equal(gotBody, audio) {
		t.Fatalf("body = %v, want file contents %v", gotBody, audio)
	}
	if out.RequestID != "req-file" {
		t.Fatalf("request_id = %q, want req-file", out.RequestID)
	}
}

func TestFromURL_AdditionalQueryParamsOverrideTypedField(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"request_id":"req-123"}`)
	}))
	defer srv.Close()

	_, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).FromURL(
		context.Background(), "https://example.invalid/audio.wav",
		&PreRecordedTranscriptionOptions{
			Model:                 "nova-3",
			AdditionalQueryParams: url.Values{"model": []string{"nova-2-meeting"}},
		})
	if err != nil {
		t.Fatalf("FromURL: %v", err)
	}
	if got := gotQuery["model"]; len(got) != 1 || got[0] != "nova-2-meeting" {
		t.Fatalf("model query = %v, want [nova-2-meeting]", got)
	}
}
