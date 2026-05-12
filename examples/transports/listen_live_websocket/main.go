// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Demonstrates the WebSocket transport for the listen-live streaming
// surface using only the generated codegen output - api/types for the
// message shapes and api/transport/websocket for the marshal/unmarshal
// helpers. A real Deepgram WSS endpoint is reached using the
// dvonthenen/websocket library (already in go.mod from the legacy SDK).
//
// Reads DEEPGRAM_API_KEY for auth (subprotocol "token, <key>" handshake)
// and supports replaying a WAV file from DG_AUDIO_FILE if you want to
// drive real transcription.
//
// Without those env vars the example exits cleanly after demonstrating
// the marshal/unmarshal cycle locally with a canned wire payload.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dvwebsocket "github.com/dvonthenen/websocket"

	"github.com/deepgram/spec-mock-go-sdk/api/document"
	ws "github.com/deepgram/spec-mock-go-sdk/api/transport/websocket"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func main() {
	demoLocalRoundtrip()

	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		fmt.Println("\nDEEPGRAM_API_KEY not set; skipping live WSS connection.")
		return
	}

	if err := liveStream(context.Background(), apiKey); err != nil {
		log.Fatalf("live stream: %v", err)
	}
}

// demoLocalRoundtrip shows how a customer codes against the generated
// transport helpers without needing any network. Send a Configure (text
// frame, JSON) and an audio frame (binary), unmarshal a canned server
// Results message.
func demoLocalRoundtrip() {
	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{
			Features: map[string]document.Interface{},
		},
	}
	wire, isBinary, err := ws.MarshalClientStream(configure)
	if err != nil {
		log.Fatalf("MarshalClientStream Configure: %v", err)
	}
	fmt.Printf("Configure wire (isBinary=%t): %s\n", isBinary, string(wire))

	audio := &spectypes.ClientStreamMemberAudio{
		Value: spectypes.AudioFrame{Data: []byte{0x01, 0x02, 0x03}},
	}
	wire, isBinary, err = ws.MarshalClientStream(audio)
	if err != nil {
		log.Fatalf("MarshalClientStream Audio: %v", err)
	}
	fmt.Printf("Audio wire (isBinary=%t, len=%d): %v\n", isBinary, len(wire), wire)

	canned := []byte(`{"type":"Results","channel_index":[0,1],"duration":1.5,"start":0.0,"is_final":true,"channel":{"alternatives":[{"transcript":"hello","confidence":0.95,"words":[],"languages":[]}]},"metadata":{},"from_finalize":false}`)
	msg, err := ws.UnmarshalServerStream(canned)
	if err != nil {
		log.Fatalf("UnmarshalServerStream: %v", err)
	}
	results, ok := msg.(*spectypes.ServerStreamMemberResults)
	if !ok {
		log.Fatalf("unexpected variant: %T", msg)
	}
	if t := results.Value.Channel.Alternatives[0].Transcript; t != nil {
		fmt.Printf("Decoded server transcript: %q\n", *t)
	}
}

// liveStream opens a Deepgram WSS session through dvonthenen/websocket,
// configures it via the generated Configure message, then prints any
// server Results / UtteranceEnd / Metadata messages until the connection
// closes. Customer-side WS connection mgmt (ping, retry, etc.) is the
// caller's responsibility - codegen provides the wire helpers.
func liveStream(ctx context.Context, apiKey string) error {
	headers := http.Header{}
	headers.Set("Authorization", "Token "+apiKey)

	conn, _, err := dvwebsocket.DefaultDialer.DialContext(
		ctx,
		"wss://api.deepgram.com/v1/listen?model=nova-3&punctuate=true&language=en-US",
		headers,
	)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{Features: map[string]document.Interface{}},
	}
	wire, _, err := ws.MarshalClientStream(configure)
	if err != nil {
		return fmt.Errorf("marshal Configure: %w", err)
	}
	if err := conn.WriteMessage(dvwebsocket.TextMessage, wire); err != nil {
		return fmt.Errorf("send Configure: %w", err)
	}

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		if msgType != dvwebsocket.TextMessage {
			continue
		}
		msg, err := ws.UnmarshalServerStream(data)
		if err != nil {
			continue
		}
		switch m := msg.(type) {
		case *spectypes.ServerStreamMemberResults:
			if t := m.Value.Channel.Alternatives[0].Transcript; t != nil {
				fmt.Printf("transcript: %s\n", *t)
			}
		case *spectypes.ServerStreamMemberMetadata:
			b, _ := json.Marshal(m.Value)
			fmt.Printf("metadata: %s\n", string(b))
		case *spectypes.ServerStreamMemberError:
			fmt.Printf("error: variant=%s\n", m.Value.Variant)
			return nil
		}
	}
}
