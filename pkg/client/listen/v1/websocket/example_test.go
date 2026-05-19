// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1_test

import (
	"fmt"

	ws "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func ExampleNewWithDefaults() {
	client := ws.NewWithDefaults()
	_ = client
	fmt.Println("constructed listen v1 websocket client")
	// Output: constructed listen v1 websocket client
}

func ExampleLiveTranscriptionOptions() {
	opts := &ws.LiveTranscriptionOptions{
		Model:          "nova-3",
		Language:       "en-US",
		Encoding:       "linear16",
		SampleRate:     16000,
		Channels:       1,
		InterimResults: true,
		SmartFormat:    true,
		VadEvents:      true,
	}
	_ = opts
	fmt.Println("ok")
	// Output: ok
}
