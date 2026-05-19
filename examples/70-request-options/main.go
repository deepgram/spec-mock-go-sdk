// Example: Request Options
//
// Every option on PreRecordedTranscriptionOptions corresponds to one
// @httpQuery on the Smithy spec. There is no arbitrary-param escape
// hatch — adding a new field requires a spec change.

package main

import (
	"context"
	"fmt"
	"log"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

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
		fmt.Println(response.Results.Channels[0].Alternatives[0].Transcript)
	}
}
