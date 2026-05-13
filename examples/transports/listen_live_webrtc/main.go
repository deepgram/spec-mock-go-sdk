// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Demonstrates the typed surface of the WebRTC transport for
// listen-live. The transport package compiles and exports the same
// Open*/Send/Recv/Close shape as the WebSocket transport, but the
// bodies return a "not implemented" error - the spec lists webrtc in
// @supportsTransports but doesn't yet model the WebRTC-specific framing
// (signaling, ICE, data channel envelope).
//
// When the spec gains WebRTC framing, spec-codegen-go fills in the
// bodies; customer code can already depend on the typed surface today.
package main

import (
	"github.com/deepgram/spec-mock-go-sdk/api/document"
	"context"
	"fmt"

	wrtc "github.com/deepgram/spec-mock-go-sdk/api/transport/webrtc"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func main() {
	ctx := context.Background()

	stream, err := wrtc.OpenStream[spectypes.ClientStream, spectypes.ServerStream](
		ctx, "wss://example.invalid/signaling",
		spectypes.MarshalClientStream,
		spectypes.UnmarshalServerStream)
	if err != nil {
		fmt.Printf("wrtc.OpenStream: %v\n", err)
		return
	}
	defer stream.Close()

	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{Features: map[string]document.Interface{}},
	}
	if err := stream.Send(configure); err != nil {
		fmt.Printf("Send: %v\n", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Printf("Recv: %v\n", err)
			return
		}
		fmt.Printf("got: %T\n", msg)
	}
}
