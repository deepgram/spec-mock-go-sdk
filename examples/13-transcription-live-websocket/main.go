// Example: Live transcription over WebSocket
package main

import (
	"context"
	"io"
	"log"
	"os"

	live "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/live"
)

func main() {
	client := live.New(live.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	stream, err := client.Connect(context.Background(), &live.LiveTranscriptionOptions{
		Model:          "nova-3",
		Encoding:       "linear16",
		SampleRate:     16000,
		InterimResults: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	go func() {
		for _, chunk := range readAudioChunks() {
			if err := stream.SendAudio(chunk); err != nil {
				return
			}
		}
		_ = stream.CloseStream()
	}()

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		switch event.(type) {
		case *live.ResultsEvent:
		case *live.MetadataEvent:
			return
		case *live.ErrorEvent:
			return
		}
	}
}

func readAudioChunks() [][]byte {
	return [][]byte{[]byte("audio bytes")}
}
