// Example: Prerecorded transcription from a file
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

	resp, err := client.FromFile(context.Background(), "examples/fixtures/audio.wav", "audio/wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model:       "nova-3",
		Punctuate:   true,
		SmartFormat: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if resp.Results != nil && len(resp.Results.Channels) > 0 && len(resp.Results.Channels[0].Alternatives) > 0 {
		fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
	}
	if resp.Metadata != nil {
		fmt.Println("request_id:", resp.Metadata.RequestID)
	}
}
