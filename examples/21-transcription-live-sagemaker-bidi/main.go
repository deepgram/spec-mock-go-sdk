// Example: Live transcription with SageMaker bidirectional streaming
//
// SageMaker endpoints typically run with network isolation enabled, so
// audio must be streamed as bytes through InvokeEndpointWithBidirectionalStream.
// Provide raw PCM (linear16, mono, 16 kHz) at AUDIO_PATH.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntimehttp2"
	live "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/live"
)

func main() {
	endpoint := os.Getenv("SAGEMAKER_ENDPOINT_NAME")
	if endpoint == "" {
		log.Fatal("SAGEMAKER_ENDPOINT_NAME is required")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}
	audioPath := os.Getenv("AUDIO_PATH")
	if audioPath == "" {
		log.Fatal("AUDIO_PATH is required (raw PCM, linear16, mono, 16 kHz)")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}
	awsClient := sagemakerruntimehttp2.NewFromConfig(cfg)
	client := live.New(live.WithSageMakerBidiTransport(awsClient, endpoint))

	stream, err := client.Connect(ctx, &live.LiveTranscriptionOptions{
		Model:          "nova-3",
		Encoding:       "linear16",
		SampleRate:     16000,
		InterimResults: true,
	})
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer stream.Close()

	f, err := os.Open(audioPath)
	if err != nil {
		log.Fatalf("open audio: %v", err)
	}
	defer f.Close()

	go func() {
		buf := make([]byte, 8000) // 250 ms @ 16 kHz linear16
		for {
			n, err := f.Read(buf)
			if n > 0 {
				if sendErr := stream.SendAudio(buf[:n]); sendErr != nil {
					log.Printf("send audio: %v", sendErr)
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("read audio: %v", err)
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
			log.Fatalf("recv: %v", err)
		}
		fmt.Printf("%T %+v\n", event, event)
		if _, done := event.(*live.MetadataEvent); done {
			return
		}
	}
}
