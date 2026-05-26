// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package livev1_test

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntimehttp2"
	live "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/live"
)

func ExampleClient_Connect() {
	client := live.New(live.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	stream, err := client.Connect(context.Background(), &live.LiveTranscriptionOptions{
		Model:          "nova-3",
		Encoding:       "linear16",
		SampleRate:     16000,
		InterimResults: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
}

func ExampleStream_SendAudio() {
	var stream *live.Stream
	if err := stream.SendAudio([]byte("audio bytes")); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_Recv() {
	var stream *live.Stream
	for {
		event, err := stream.Recv()
		if err != nil {
			return
		}
		switch event.(type) {
		case *live.ResultsEvent:
		case *live.MetadataEvent:
			return
		case *live.ErrorEvent:
			return
		}
	}
}

func ExampleStream_Finalize() {
	var stream *live.Stream
	if err := stream.Finalize(-1); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_KeepAlive() {
	var stream *live.Stream
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := stream.KeepAlive(); err != nil {
			return
		}
	}
}

func ExampleStream_CloseStream() {
	var stream *live.Stream
	if err := stream.CloseStream(); err != nil {
		log.Fatal(err)
	}
}

func ExampleLiveTranscriptionOptions() {
	opts := &live.LiveTranscriptionOptions{
		Model:                 "nova-3",
		Language:              "en-US",
		Encoding:              "linear16",
		SampleRate:            16000,
		Channels:              1,
		InterimResults:        true,
		SmartFormat:           true,
		Punctuate:             true,
		VadEvents:             true,
		UtteranceEndMs:        1000,
		AdditionalQueryParams: url.Values{"experimental": []string{"true"}},
	}
	_ = opts
}

func ExampleWithSageMakerBidiTransport() {
	awsClient := sagemakerruntimehttp2.New(sagemakerruntimehttp2.Options{Region: "us-east-1"})
	client := live.New(live.WithSageMakerBidiTransport(awsClient, "endpoint-name"))
	stream, err := client.Connect(context.Background(), &live.LiveTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
}
