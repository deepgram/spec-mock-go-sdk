// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package prerecordedv1_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func ExampleClient_FromURL() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	resp, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model:       "nova-3",
		Punctuate:   true,
		SmartFormat: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.RequestID)
}

func ExampleClient_FromFile() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	resp, err := client.FromFile(context.Background(), "examples/fixtures/audio.wav", "audio/wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model: "nova-3",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.RequestID)
}

func ExampleClient_FromStream() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	audio := strings.NewReader("audio bytes")
	resp, err := client.FromStream(context.Background(), audio, "audio/wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model: "nova-3",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.RequestID)
}

func ExamplePreRecordedTranscriptionOptions() {
	opts := &prerecorded.PreRecordedTranscriptionOptions{
		Model:                 "nova-3",
		Language:              "en-US",
		Punctuate:             true,
		SmartFormat:           true,
		Diarize:               true,
		Utterances:            true,
		Paragraphs:            true,
		DetectEntities:        true,
		Sentiment:             true,
		Tag:                   []string{"demo"},
		AdditionalQueryParams: url.Values{"experimental": []string{"true"}},
	}
	_ = opts
}

func ExampleWithSageMakerTransport() {
	awsClient := sagemakerruntime.New(sagemakerruntime.Options{Region: "us-east-1"})
	client := prerecorded.New(prerecorded.WithSageMakerTransport(awsClient, "endpoint-name", prerecorded.WithTargetVariant("variantA")))
	resp, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model: "nova-3",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.RequestID)
}
