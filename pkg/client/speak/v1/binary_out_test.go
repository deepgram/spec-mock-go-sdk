// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the binary-audio-out transport path: the Speak
// response body is raw audio bytes, with metadata in dg-* response headers.
// Verifies Invoke reads the body into Audio and binds the @httpHeader fields,
// rather than JSON-decoding (which would fail on binary).
package speakv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBinaryAudioOut(t *testing.T) {
	audio := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x01, 0x02, 0x03} // fake audio bytes

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("dg-model-name", "aura-2-asteria-en")
		w.Header().Set("dg-char-count", "5")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(audio)
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL), WithAPIKey("test-key"))
	resp, err := c.FromText(context.Background(), "hello", &SpeakOptions{Model: "aura-2-asteria-en"})
	if err != nil {
		t.Fatalf("FromText returned error: %v", err)
	}
	if string(resp.Audio) != string(audio) {
		t.Fatalf("audio body not read into Audio: got %v want %v", resp.Audio, audio)
	}
	if resp.ContentType == nil || *resp.ContentType != "audio/mpeg" {
		t.Fatalf("Content-Type header not bound: %v", resp.ContentType)
	}
	if resp.ModelName == nil || *resp.ModelName != "aura-2-asteria-en" {
		t.Fatalf("dg-model-name header not bound: %v", resp.ModelName)
	}
	if resp.CharCount == nil || *resp.CharCount != 5 {
		t.Fatalf("dg-char-count header not bound/parsed: %v", resp.CharCount)
	}
}
