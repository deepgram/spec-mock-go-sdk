// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	ws "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func ExampleNew() {
	client := ws.New("your-api-key", "")
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleNewWithDefaults() {
	_ = os.Setenv("DEEPGRAM_API_KEY", "your-api-key")
	client := ws.NewWithDefaults()
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleClient_WithBaseURL() {
	client := ws.NewWithDefaults().
		WithBaseURL("wss://staging.api.deepgram.com")
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleClient_Connect() {
	client := ws.NewWithDefaults()
	stream, err := client.Connect(context.Background(),
		&ws.LiveTranscriptionOptions{
			Model:          "nova-3",
			Encoding:       "linear16",
			SampleRate:     16000,
			InterimResults: true,
			SmartFormat:    true,
		})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
}

func ExampleStream_SendAudio() {
	var stream *ws.Stream
	chunk := make([]byte, 4096)
	if err := stream.SendAudio(chunk); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_CloseStream() {
	var stream *ws.Stream
	if err := stream.CloseStream(); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_Finalize() {
	var stream *ws.Stream
	if err := stream.Finalize(-1); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_KeepAlive() {
	var stream *ws.Stream
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := stream.KeepAlive(); err != nil {
			return
		}
	}
}

func ExampleStream_Sync() {
	var stream *ws.Stream
	if err := stream.Sync(42); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_Recv() {
	var stream *ws.Stream
	for {
		event, err := stream.Recv()
		if err != nil {
			return
		}
		switch e := event.(type) {
		case *ws.ResultsEvent:
			if e.IsFinal {
				fmt.Println(e.Channel.Alternatives[0].Transcript)
			}
		case *ws.MetadataEvent:
			return
		case *ws.ErrorEvent:
			log.Fatal(e.Description)
		case *ws.SpeechStartedEvent:
			fmt.Println("speech started at", e.Timestamp)
		case *ws.UtteranceEndEvent:
			fmt.Println("utterance ended at", e.LastWordEnd)
		}
	}
}

func ExampleStream_Close() {
	var stream *ws.Stream
	if err := stream.Close(); err != nil {
		log.Fatal(err)
	}
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
		Punctuate:      true,
		VadEvents:      true,
		UtteranceEndMs: 1000,
		Endpointing:    500,
		Diarize:        true,
		Tag:            []string{"my-session"},
		Keyterm:        []string{"deepgram"},
	}
	_ = opts
	fmt.Println("ok")
	// Output: ok
}
