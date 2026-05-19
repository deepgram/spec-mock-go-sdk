// Example: Transcribe Prerecorded Audio with Callback
//
// Setting Callback turns the request into an async fire-and-forget
// dispatch: /v1/listen returns immediately with a request_id, then
// Deepgram POSTs the actual transcription to the supplied URL when
// processing completes.
//
// You need a publicly reachable webhook to receive the result.
// Replace the placeholder URL below before running.

package main

import (
	"context"
	"fmt"
	"log"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	fmt.Println("Sending transcription request with callback...")
	response, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model:    "nova-3",
			Callback: "https://your-callback-url.com/webhook",
		},
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Request accepted. Request ID:", response.Metadata.RequestID)
	fmt.Println("Transcription will be POSTed to your callback URL when ready.")
}
