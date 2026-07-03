// Command speak-smoke is an end-to-end smoke test of the SpeakLive (TTS) facade
// over the bundled SageMaker transport: it constructs the real speak/v1/live
// client with WithSageMakerBidiTransport, sends text, and confirms synthesized
// audio frames come back — exercising the facade -> ResilientDialBidi -> live
// SageMaker path for the binary-on-server (TTS) product.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	livev1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/speak/v1/live"
)

func main() {
	endpoint := flag.String("endpoint", "dg-tts-test", "SageMaker endpoint")
	region := flag.String("region", "us-east-2", "AWS region")
	model := flag.String("model", "aura-2-thalia-en", "Aura model")
	text := flag.String("text", "Hello world. This is a self hosted text to speech smoke test.", "text to synthesize")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	if err != nil {
		fmt.Fprintln(os.Stderr, "aws config:", err)
		os.Exit(1)
	}

	client := livev1.New(livev1.WithSageMakerBidiTransport(cfg, *endpoint))
	stream, err := client.Connect(ctx, &livev1.SpeakLiveOptions{Model: *model, Encoding: "linear16", SampleRate: 16000})
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect:", err)
		os.Exit(1)
	}
	defer stream.Close()

	var audioBytes, frames int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			ev, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					fmt.Fprintln(os.Stderr, "recv:", err)
				}
				return
			}
			if a, ok := ev.(*livev1.AudioEvent); ok {
				frames++
				audioBytes += len(a.Data)
			}
		}
	}()

	if err := stream.SendSpeak(*text); err != nil {
		fmt.Fprintln(os.Stderr, "send:", err)
		os.Exit(1)
	}
	_ = stream.SendFlush()
	// Give the model time to synthesize + stream the audio back.
	select {
	case <-done:
	case <-time.After(10 * time.Second):
	}
	_ = stream.SendClose()

	fmt.Printf("speak-smoke: audioFrames=%d audioBytes=%d\n", frames, audioBytes)
	if audioBytes == 0 {
		os.Exit(1)
	}
}
