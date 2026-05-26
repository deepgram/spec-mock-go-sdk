// Example: Authentication with an API key
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
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
