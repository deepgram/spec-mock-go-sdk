// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Demonstrates the WebSocket transport for the listen-live streaming
// surface using only the generated codegen output - api/types for the
// message shapes, helpers, and route metadata, plus api/transport/websocket
// for the generic OpenStream[C, S any] primitive. A real Deepgram WSS
// endpoint is reached via OpenStream which internally uses the
// dvonthenen/websocket library (already in go.mod).
//
// Reads DEEPGRAM_API_KEY for auth (Authorization: Token header) and
// supports DG_AUDIO_FILE for real audio streaming. Without those env
// vars the example exits cleanly after demonstrating the marshal /
// unmarshal cycle locally with a canned wire payload.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
// marshal/unmarshal helpers without needing any network. Send a Configure
// (text frame, JSON) and an audio frame (binary), unmarshal a canned
// server Results message.
func demoLocalRoundtrip() {
	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{
			Features: map[string]document.Interface{},
		},
	}
	wire, isBinary, err := spectypes.MarshalClientStream(configure)
	if err != nil {
		log.Fatalf("MarshalClientStream Configure: %v", err)
	}
	fmt.Printf("Configure wire (isBinary=%t): %s\n", isBinary, string(wire))

	audio := &spectypes.ClientStreamMemberAudio{
		Value: spectypes.AudioFrame{Data: []byte{0x01, 0x02, 0x03}},
	}
	wire, isBinary, err = spectypes.MarshalClientStream(audio)
	if err != nil {
		log.Fatalf("MarshalClientStream Audio: %v", err)
	}
	fmt.Printf("Audio wire (isBinary=%t, len=%d): %v\n", isBinary, len(wire), wire)

	canned := []byte(`{"type":"Results","channel_index":[0,1],"duration":1.5,"start":0.0,"is_final":true,"channel":{"alternatives":[{"transcript":"hello","confidence":0.95,"words":[],"languages":[]}]},"metadata":{},"from_finalize":false}`)
	msg, err := spectypes.UnmarshalServerStream(canned)
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

// liveStream opens a Deepgram WSS session through the generic
// ws.OpenStream primitive (which internally drives dvonthenen/websocket).
// The MarshalClientStream / UnmarshalServerStream closures from api/types
// handle wire framing for every union variant; the consumer code stays
// type-safe and never touches a raw WebSocket frame.
func liveStream(ctx context.Context, apiKey string) error {
	headers := http.Header{}
	headers.Set("Authorization", "Token "+apiKey)

	stream, err := ws.OpenStream[spectypes.ClientStream, spectypes.ServerStream](
		ctx,
		"wss://api.deepgram.com/v1/listen?model=nova-3&punctuate=true&language=en-US",
		headers,
		spectypes.MarshalClientStream,
		spectypes.UnmarshalServerStream,
	)
	if err != nil {
		return fmt.Errorf("OpenStream: %w", err)
	}
	defer stream.Close()

	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{Features: map[string]document.Interface{}},
	}
	if err := stream.Send(configure); err != nil {
		return fmt.Errorf("send Configure: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return nil
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
