// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Runtime safety-net for the Agent's binary-on-BOTH-sides path: the client
// streams mic audio (binary in) + control JSON; the server streams agent
// speech (binary out) + event JSON. Verifies SendAudio/SendSettings reach the
// server, and that Recv routes a server binary frame to the audio union member
// and a JSON frame to the Welcome member.
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

func TestAgentBinaryBothWays(t *testing.T) {
	agentSpeech := []byte{0xAA, 0xBB, 0xCC}
	upgrader := ws.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		// Client sends: Settings (text), then mic audio (binary).
		mt1, _, err := c.ReadMessage()
		if err != nil || mt1 != ws.TextMessage {
			return
		}
		mt2, _, err := c.ReadMessage()
		if err != nil || mt2 != ws.BinaryMessage {
			return
		}
		// Server replies: Welcome (text), then agent speech (binary).
		_ = c.WriteMessage(ws.TextMessage, []byte(`{"type":"Welcome","request_id":"00000000-0000-0000-0000-000000000000"}`))
		_ = c.WriteMessage(ws.BinaryMessage, agentSpeech)
	}))
	defer srv.Close()

	baseURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	stream, err := New(WithBaseURL(baseURL), WithAPIKey("k")).
		Connect(context.Background(), &AgentLiveOptions{})
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	if err := stream.SendSettings(spectypes.AgentSettings{}); err != nil {
		t.Fatalf("SendSettings: %v", err)
	}
	if err := stream.SendAudio([]byte{0x01, 0x02}); err != nil { // mic audio in (binary)
		t.Fatalf("SendAudio: %v", err)
	}

	m1, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (welcome): %v", err)
	}
	if _, ok := m1.(*WelcomeEvent); !ok {
		t.Fatalf("expected Welcome member, got %T", m1)
	}

	m2, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv (audio): %v", err)
	}
	af, ok := m2.(*AudioEvent)
	if !ok {
		t.Fatalf("expected binary audio member, got %T", m2)
	}
	if string(af.Data) != string(agentSpeech) {
		t.Fatalf("agent speech mismatch: got %v want %v", af.Data, agentSpeech)
	}
}
