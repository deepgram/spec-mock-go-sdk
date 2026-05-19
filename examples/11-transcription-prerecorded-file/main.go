// Example: Transcribe Prerecorded Audio from File
//
// FromFile reads the file from disk and streams it to /v1/listen as
// the request body. For in-memory bytes or non-file readers, use
// FromStream with an explicit content type.

package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	_, thisFile, _, _ := runtime.Caller(0)
	audioPath := filepath.Join(filepath.Dir(thisFile), "..", "fixtures", "audio.wav")

	fmt.Println("Reading audio file:", audioPath)
	fmt.Println("Sending transcription request...")

	response, err := client.FromFile(
		context.Background(),
		audioPath,
		"audio/wav",
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
