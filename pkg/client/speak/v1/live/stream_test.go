// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the live (WSS) binary-on-server path: the SpeakLive
// server streams synthesized audio as binary frames interleaved with JSON
// status messages. Verifies the facade sends a Speak message, and that Recv
// routes a binary frame to the audio union member and a text frame to the JSON
// member — exercising the transport's isBinary-aware unmarshal.
package livev1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	ws "github.com/gorilla/websocket"
)

func TestSpeakLiveBinaryRecv(t *testing.T) {
	audio := []byte{0x01, 0x02, 0x03, 0x04}
	upgrader := ws.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		if _, _, err := c.ReadMessage(); err != nil { // the client's Speak message
			return
		}
		_ = c.WriteMessage(ws.BinaryMessage, audio)                                  // audio out
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"Flushed","sequence_id":1}`)) // status
	}))
	defer srv.Close()

	baseURL := "ws" + strings.TrimPrefix(srv.URL, "http") // ws://127.0.0.1:port
	stream, err := New(WithBaseURL(baseURL), WithAPIKey("k")).
		Connect(context.Background(), &SpeakLiveOptions{Model: "aura-2-asteria-en"})
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	if err := stream.SendSpeak("hello"); err != nil {
		t.Fatalf("SendSpeak: %v", err)
	}

	m1, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (audio): %v", err)
	}
	af, ok := m1.(*spectypes.SpeakServerStreamMemberAudio)
	if !ok {
		t.Fatalf("expected binary audio member, got %T", m1)
	}
	if string(af.Value.Data) != string(audio) {
		t.Fatalf("audio bytes mismatch: got %v want %v", af.Value.Data, audio)
	}

	m2, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (status): %v", err)
	}
	if _, ok := m2.(*spectypes.SpeakServerStreamMemberFlushed); !ok {
		t.Fatalf("expected Flushed JSON member, got %T", m2)
	}
}
