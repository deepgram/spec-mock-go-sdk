// Example: Live transcription with SageMaker bidirectional streaming
package main

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntimehttp2"
	live "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/live"
)

func main() {
	awsClient := sagemakerruntimehttp2.New(sagemakerruntimehttp2.Options{Region: "us-east-1"})
	client := live.New(live.WithSageMakerBidiTransport(awsClient, "endpoint-name"))

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
		_ = stream.SendAudio([]byte("audio bytes"))
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
		if _, done := event.(*live.MetadataEvent); done {
			return
		}
	}
}
