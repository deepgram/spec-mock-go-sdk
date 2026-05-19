// Example: Transcription with Advanced Options
//
// Demonstrates the option fields most production integrations want:
// smart formatting (numbers, dates, addresses), punctuation,
// diarization (per-word speaker labels), explicit language pin.
//
// The full inventory of options lives in
// pkg/client/listen/v1/rest/options.go. Each field maps to one
// @httpQuery member on the Smithy spec's TranscribeInput.

package main

import (
	"context"
	"fmt"
	"log"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	fmt.Println("Transcribing with advanced options...")
	response, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model:       "nova-3",
			Language:    "en-US",
			SmartFormat: true,
			Punctuate:   true,
			Diarize:     true,
		},
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if response.Results == nil || len(response.Results.Channels) == 0 {
		fmt.Println("No results returned")
		return
	}
	channel := response.Results.Channels[0]
	if len(channel.Alternatives) == 0 {
		return
	}
	alt := channel.Alternatives[0]
	fmt.Printf("Transcript: %s\n", alt.Transcript)

	if len(alt.Words) > 0 {
		fmt.Println("\nSpeaker diarization:")
		for _, word := range alt.Words {
			if word.Speaker != nil {
				fmt.Printf("  Speaker %d: %s\n", *word.Speaker, word.Word)
			}
		}
	}
}
