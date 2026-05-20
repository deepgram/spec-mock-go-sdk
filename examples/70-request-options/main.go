// Example: Request Options — arbitrary query parameters.
//
// PreRecordedTranscriptionOptions.AdditionalQueryParams is the escape
// hatch for sending query parameters the SDK does not yet expose as
// typed fields. Use it when Deepgram ships a new parameter on the
// API before the SDK has been updated to recognise it. On collision
// with a typed field, AdditionalQueryParams wins.

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	response, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model:       "nova-3",
			SmartFormat: true,
			AdditionalQueryParams: url.Values{
				"experimental_feature": []string{"true"},
				"custom_tag":           []string{"a", "b"},
			},
		},
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if response.Results == nil || len(response.Results.Channels) == 0 {
		return
	}
	if len(response.Results.Channels[0].Alternatives) > 0 {
		fmt.Println(response.Results.Channels[0].Alternatives[0].Transcript)
	}
}
