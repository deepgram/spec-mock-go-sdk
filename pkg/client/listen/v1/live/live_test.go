// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Listen live (WebSocket) facade: the client opens
// a WS upgrade carrying the typed options as query params, streams binary
// audio frames, and receives JSON server messages. Verifies the upgrade query
// wiring, that SendAudio reaches the server, and that Recv routes the JSON
// Results + Metadata frames to their sealed Event variants.
package livev1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ws "github.com/gorilla/websocket"
)

func TestLiveStreamWireAndEvents(t *testing.T) {
	var gotQuery map[string][]string
	var gotAudioType int
	upgrader := ws.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		mt, _, err := c.ReadMessage() // client mic audio (binary frame)
		if err != nil {
			return
		}
		gotAudioType = mt
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"Results","is_final":true}`))
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"Metadata","request_id":"req-meta"}`))
	}))
	defer srv.Close()

	baseURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	stream, err := New(WithBaseURL(baseURL), WithAPIKey("k")).Connect(
		context.Background(), &LiveTranscriptionOptions{
			Model:          "nova-3",
			Encoding:       "linear16",
			SampleRate:     16000,
			InterimResults: true,
		})
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	if err := stream.SendAudio([]byte{0x01, 0x02, 0x03}); err != nil {
		t.Fatalf("SendAudio: %v", err)
	}

	m1, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (results): %v", err)
	}
	res, ok := m1.(*ResultsEvent)
	if !ok {
		t.Fatalf("expected *ResultsEvent, got %T", m1)
	}
	if res.IsFinal == nil || !*res.IsFinal {
		t.Fatalf("is_final = %v, want true", res.IsFinal)
	}

	m2, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (metadata): %v", err)
	}
	meta, ok := m2.(*MetadataEvent)
	if !ok {
		t.Fatalf("expected *MetadataEvent, got %T", m2)
	}
	if meta.RequestId == nil || *meta.RequestId != "req-meta" {
		t.Fatalf("request_id = %v, want req-meta", meta.RequestId)
	}

	// The audio frame must reach the server as a binary WS frame.
	if gotAudioType != ws.BinaryMessage {
		t.Fatalf("audio frame type = %d, want binary (%d)", gotAudioType, ws.BinaryMessage)
	}
	// Typed options must ride on the upgrade URL query.
	for name, want := range map[string]string{
		"model": "nova-3", "encoding": "linear16",
		"sample_rate": "16000", "interim_results": "true",
	} {
		if q := gotQuery[name]; len(q) != 1 || q[0] != want {
			t.Fatalf("%s query = %v, want [%s]", name, q, want)
		}
	}

	if err := stream.CloseStream(); err != nil {
		t.Fatalf("CloseStream: %v", err)
	}
}
