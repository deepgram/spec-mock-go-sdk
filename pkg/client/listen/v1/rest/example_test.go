// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	rest "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func ExampleNew() {
	client := rest.New("your-api-key", "")
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleNewWithDefaults() {
	_ = os.Setenv("DEEPGRAM_API_KEY", "your-api-key")
	client := rest.NewWithDefaults()
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleClient_WithBaseURL() {
	client := rest.NewWithDefaults().
		WithBaseURL("https://staging.api.deepgram.com")
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleClient_WithHTTPClient() {
	hc := &http.Client{Timeout: 30 * time.Second}
	client := rest.NewWithDefaults().WithHTTPClient(hc)
	_ = client
	fmt.Println("ok")
	// Output: ok
}

func ExampleClient_FromURL() {
	client := rest.NewWithDefaults()
	resp, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&rest.PreRecordedTranscriptionOptions{
			Model:       "nova-3",
			Punctuate:   true,
			SmartFormat: true,
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
}

func ExampleClient_FromFile() {
	client := rest.NewWithDefaults()
	resp, err := client.FromFile(context.Background(),
		"./recording.wav", "audio/wav",
		&rest.PreRecordedTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
}

func ExampleClient_FromStream() {
	client := rest.NewWithDefaults()
	audio := strings.NewReader("...audio bytes...")
	resp, err := client.FromStream(context.Background(),
		audio, "audio/wav",
		&rest.PreRecordedTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
}

func ExamplePreRecordedTranscriptionOptions() {
	opts := &rest.PreRecordedTranscriptionOptions{
		Model:          "nova-3",
		Language:       "en-US",
		Punctuate:      true,
		SmartFormat:    true,
		Utterances:     true,
		Paragraphs:     true,
		Diarize:        true,
		FillerWords:    false,
		DetectEntities: true,
		Sentiment:      true,
		Topics:         true,
		Intents:        true,
		Summarize:      "v2",
		Tag:            []string{"my-tag"},
		Keyterm:        []string{"deepgram"},
		Redact:         []string{"pii"},
	}
	_ = opts
	fmt.Println("ok")
	// Output: ok
}
