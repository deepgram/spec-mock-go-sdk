// Example: Advanced transcription options
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	opts := &prerecorded.PreRecordedTranscriptionOptions{
		Model:          "nova-3",
		Language:       "en-US",
		Punctuate:      true,
		SmartFormat:    true,
		Diarize:        true,
		Utterances:     true,
		Paragraphs:     true,
		DetectEntities: true,
		Sentiment:      true,
		Topics:         true,
		Intents:        true,
		Tag:            []string{"demo"},
		AdditionalQueryParams: url.Values{
			"experimental": []string{"true"},
		},
	}

	resp, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.RequestID)
}
