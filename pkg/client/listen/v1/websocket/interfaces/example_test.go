// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfacesv1_test

import (
	"fmt"

	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket/interfaces"
)

// ExampleMessageResponse shows the customer-facing transcription
// event surfaced from a Listen v1 WebSocket stream. Handlers receive
// this via LiveMessageChan or LiveMessageCallback. Read transcripts
// off the first Alternative on the embedded Channel.
func ExampleMessageResponse() {
	resp := &interfaces.MessageResponse{
		Type:    "Results",
		IsFinal: true,
		Channel: interfaces.Channel{
			Alternatives: []interfaces.Alternative{
				{
					Transcript: "hello world",
					Confidence: 0.98,
				},
			},
		},
	}
	if len(resp.Channel.Alternatives) > 0 {
		alt := resp.Channel.Alternatives[0]
		fmt.Printf("final=%t transcript=%q confidence=%.2f\n",
			resp.IsFinal, alt.Transcript, alt.Confidence)
	}
	// Output: final=true transcript="hello world" confidence=0.98
}

// ExampleMetadataResponse demonstrates the per-session metadata event
// emitted at session end. RequestID correlates with the dg-request-id
// header on the upstream handshake.
func ExampleMetadataResponse() {
	m := &interfaces.MetadataResponse{
		Type:      "Metadata",
		RequestID: "req-ws-456",
		Channels:  1,
		Duration:  4.2,
	}
	fmt.Printf("type=%s request=%s duration=%.1fs channels=%d\n",
		m.Type, m.RequestID, m.Duration, m.Channels)
	// Output: type=Metadata request=req-ws-456 duration=4.2s channels=1
}

// ExampleUtteranceEndResponse shows a VAD-derived utterance boundary
// event. Customer code uses LastWordEnd to align rendering or
// downstream processing with the end of an utterance.
func ExampleUtteranceEndResponse() {
	u := &interfaces.UtteranceEndResponse{
		Type:        "UtteranceEnd",
		Channel:     []int{0},
		LastWordEnd: 3.14,
	}
	fmt.Printf("type=%s last_word_end=%.2f\n", u.Type, u.LastWordEnd)
	// Output: type=UtteranceEnd last_word_end=3.14
}
