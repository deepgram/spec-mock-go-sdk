// Command flux-smoke is an end-to-end smoke test of the Flux (Listen v2) facade
// over the bundled SageMaker transport: it constructs the real flux/v2/live
// client with WithSageMakerBidiTransport, streams a WAV, and confirms typed
// turn events come back — exercising the facade -> ResilientDialBidi -> live
// SageMaker path for a product other than the STT one the stress harness used.
package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	livev2 "github.com/deepgram/spec-mock-go-sdk/pkg/client/flux/v2/live"
)

func main() {
	endpoint := flag.String("endpoint", "dg-flux-test", "SageMaker endpoint")
	wavPath := flag.String("file", "", "16 kHz mono s16le WAV")
	region := flag.String("region", "us-east-2", "AWS region")
	model := flag.String("model", "flux-general-en", "Flux model")
	flag.Parse()

	pcm, sr, err := readWAV(*wavPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wav:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	if err != nil {
		fmt.Fprintln(os.Stderr, "aws config:", err)
		os.Exit(1)
	}

	client := livev2.New(livev2.WithSageMakerBidiTransport(cfg, *endpoint))
	stream, err := client.Connect(ctx, &livev2.FluxLiveOptions{Model: *model, Encoding: "linear16", SampleRate: sr})
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect:", err)
		os.Exit(1)
	}
	defer stream.Close()

	var connected, turns int
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
			switch e := ev.(type) {
			case *livev2.ConnectedEvent:
				connected++
			case *livev2.TurnInfoEvent:
				turns++
				if e.Transcript != nil && *e.Transcript != "" {
					fmt.Printf("  [TurnInfo] %q\n", *e.Transcript)
				}
			}
		}
	}()

	const chunk = 8192
	dur := time.Duration(float64(chunk/2) / float64(sr) * float64(time.Second))
	for off := 0; off < len(pcm); off += chunk {
		end := off + chunk
		if end > len(pcm) {
			end = len(pcm)
		}
		if err := stream.SendAudio(pcm[off:end]); err != nil {
			fmt.Fprintln(os.Stderr, "send:", err)
			break
		}
		time.Sleep(dur)
	}
	_ = stream.SendCloseStream()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	fmt.Printf("flux-smoke: connected=%d turnInfo=%d\n", connected, turns)
	if turns == 0 {
		os.Exit(1)
	}
}

func readWAV(path string) ([]byte, int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return nil, 0, fmt.Errorf("not a RIFF/WAVE file")
	}
	var sr int
	var pcm []byte
	for off := 12; off+8 <= len(b); {
		id := string(b[off : off+4])
		size := int(binary.LittleEndian.Uint32(b[off+4 : off+8]))
		body := off + 8
		if body+size > len(b) {
			size = len(b) - body
		}
		switch id {
		case "fmt ":
			if size >= 16 {
				sr = int(binary.LittleEndian.Uint32(b[body+4 : body+8]))
			}
		case "data":
			pcm = b[body : body+size]
		}
		off = body + size
		if size%2 == 1 {
			off++
		}
	}
	if sr == 0 || pcm == nil {
		return nil, 0, fmt.Errorf("missing fmt/data chunk")
	}
	return pcm, sr, nil
}
