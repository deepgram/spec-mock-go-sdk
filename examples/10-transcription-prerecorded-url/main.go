// Example: Transcribe Prerecorded Audio from URL
//
// Synchronously transcribe audio at a public HTTPS URL. The full
// response is returned on the same call — no callback, no polling.
// Best for short pre-recorded clips where blocking until done is fine.
//
// For long files or async fire-and-forget, see
// 12-transcription-prerecorded-callback.

package main

import (
	"context"
	"fmt"
	"log"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	fmt.Println("Sending transcription request...")
	response, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model: "nova-3",
		},
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Transcription received:")
	if response.Results != nil && len(response.Results.Channels) > 0 {
		channel := response.Results.Channels[0]
		if len(channel.Alternatives) > 0 {
			fmt.Println(channel.Alternatives[0].Transcript)
		}
	}
}
