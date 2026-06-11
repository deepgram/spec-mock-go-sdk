// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for Flux (Listen v2): client streams binary audio, server
// emits JSON turn events. Verifies SendAudio reaches the server and Recv routes
// the JSON Connected + TurnInfo events to their union members.
package livev2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	ws "github.com/gorilla/websocket"
)

func TestFluxTurnEvents(t *testing.T) {
	upgrader := ws.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		mt, _, err := c.ReadMessage() // client mic audio (binary)
		if err != nil || mt != ws.BinaryMessage {
			return
		}
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"Connected","request_id":"00000000-0000-0000-0000-000000000000","sequence_id":0}`))
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"TurnInfo","request_id":"00000000-0000-0000-0000-000000000000","event":"EndOfTurn","turn_index":0,"audio_window_start":0,"audio_window_end":1.5,"transcript":"hello world","words":[{"word":"hello","confidence":0.99}],"end_of_turn_confidence":0.95,"sequence_id":1}`))
	}))
	defer srv.Close()

	baseURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	stream, err := New(WithBaseURL(baseURL), WithAPIKey("k")).
		Connect(context.Background(), &FluxLiveOptions{Model: "flux-general-en"})
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	if err := stream.SendAudio([]byte{0x01, 0x02, 0x03}); err != nil {
		t.Fatalf("SendAudio: %v", err)
	}

	m1, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (connected): %v", err)
	}
	if _, ok := m1.(*ConnectedEvent); !ok {
		t.Fatalf("expected Connected, got %T", m1)
	}

	m2, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (turninfo): %v", err)
	}
	ti, ok := m2.(*TurnInfoEvent)
	if !ok {
		t.Fatalf("expected TurnInfo, got %T", m2)
	}
	if ti.Transcript == nil || *ti.Transcript != "hello world" {
		t.Fatalf("transcript = %v", ti.Transcript)
	}
	if ti.Event != spectypes.FluxEventEndOfTurn {
		t.Fatalf("event = %v", ti.Event)
	}
}
