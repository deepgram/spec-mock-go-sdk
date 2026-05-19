// Example: Request Options
//
// The Python SDK exposes an `additional_query_parameters` dict for
// passing arbitrary query parameters Deepgram might accept that the
// SDK does not have a typed field for. This Go SDK is spec-driven
// (every option is a typed field corresponding to a Smithy @httpQuery
// member) and intentionally does NOT have an arbitrary-param escape
// hatch — adding a new field is a spec change followed by a codegen
// roll-forward.
//
// What this example demonstrates instead: the most common reason to
// reach for arbitrary params in the Python SDK is DetectLanguage,
// which is a typed field on PreRecordedTranscriptionOptions.

package main

import (
	"context"
	"fmt"
	"log"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	fmt.Println("Transcribing with DetectLanguage (the typed equivalent of Python's request_options escape hatch)...")
	response, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model:          "nova-3",
			DetectLanguage: []string{"en", "es"},
		},
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if response.Results == nil || len(response.Results.Channels) == 0 {
		return
	}
	if len(response.Results.Channels[0].Alternatives) > 0 {
		fmt.Printf("Transcript: %s\n", response.Results.Channels[0].Alternatives[0].Transcript)
	}
}
