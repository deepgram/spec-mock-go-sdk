// Command sagemaker-stress is a load harness for the bundled SageMaker bidi
// transport (api/transport/sagemaker). It mirrors the Python reference harness
// (dg-sagemaker/python-stt/stt_wav_stress.py): N concurrent bidirectional
// streams against a SageMaker STT endpoint, streaming a 16 kHz mono PCM WAV in
// real-time-paced 8192-byte chunks, then reporting how many connections
// streamed cleanly vs errored. Used to compare our Go transport's
// connection-establishment behavior under burst against the known-good Python
// client (see divergence-log SAGEMAKER-002).
//
// Not part of the generated SDK; a dev/validation tool.
package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntimehttp2"
	sm "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
)

const chunkSize = 8192

func main() {
	endpoint := flag.String("endpoint", "", "SageMaker endpoint name")
	conns := flag.Int("connections", 1, "concurrent connections")
	wavPath := flag.String("file", "", "path to a 16 kHz mono s16le WAV")
	region := flag.String("region", "us-east-2", "AWS region")
	model := flag.String("model", "nova-3", "Deepgram model")
	lang := flag.String("language", "en", "language code")
	isolate := flag.Bool("isolate", false, "give each stream its own HTTP client (own connection pool) — the conn-per-stream / maxStreams=1 equivalent")
	flag.Parse()
	if *endpoint == "" || *wavPath == "" {
		fmt.Fprintln(os.Stderr, "usage: sagemaker-stress -endpoint NAME -file WAV [-connections N]")
		os.Exit(2)
	}

	pcm, sampleRate, err := readWAV(*wavPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wav:", err)
		os.Exit(1)
	}
	query := fmt.Sprintf("model=%s&language=%s&encoding=linear16&sample_rate=%d", *model, *lang, sampleRate)

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	if err != nil {
		fmt.Fprintln(os.Stderr, "aws config:", err)
		os.Exit(1)
	}
	marshal := func(b []byte) ([]byte, bool, error) { return b, true, nil }
	unmarshal := func(b []byte, _ bool) ([]byte, error) { return b, nil }

	// Default shared client (no Tier-B fix) — the control case that multiplexes.
	shared := sagemakerruntimehttp2.NewFromConfig(awsCfg)
	// open returns one bidi stream. With -isolate it uses the SHIPPED path
	// (sm.DialBidi: fresh isolated client per stream + connect retry); else the
	// shared multiplexing client.
	open := func() (sm.Stream[[]byte, []byte], error) {
		if *isolate {
			return sm.DialBidi[[]byte, []byte](ctx, awsCfg, sm.DefaultConfig(), *endpoint, "v1/listen", query, "", "", "", "", marshal, unmarshal)
		}
		return sm.OpenStream[[]byte, []byte](ctx, shared, *endpoint, "v1/listen", query, "", "", "", "", marshal, unmarshal)
	}

	fmt.Printf("Go transport stress: endpoint=%s connections=%d isolate=%v wav=%s (%d Hz)\n",
		*endpoint, *conns, *isolate, *wavPath, sampleRate)

	var ok, errd int64
	var mu sync.Mutex
	errCounts := map[string]int{}
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < *conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if e := runOne(open, pcm, sampleRate); e != nil {
				atomic.AddInt64(&errd, 1)
				mu.Lock()
				errCounts[shorten(e.Error())]++
				mu.Unlock()
			} else {
				atomic.AddInt64(&ok, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Printf("\n=== Go transport stress (api/transport/sagemaker) ===\n")
	fmt.Printf("connections=%d  successful=%d  errored=%d  wall=%.1fs\n",
		*conns, ok, errd, time.Since(start).Seconds())
	for msg, n := range errCounts {
		fmt.Printf("  (x%d) %s\n", n, msg)
	}
}

func runOne(open func() (sm.Stream[[]byte, []byte], error), pcm []byte, sampleRate int) error {
	stream, err := open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer stream.Close()

	recvErr := make(chan error, 1)
	go func() {
		for {
			if _, e := stream.Recv(); e != nil {
				recvErr <- e
				return
			}
		}
	}()

	framesPerChunk := chunkSize / 2 // 16-bit mono
	chunkDur := time.Duration(float64(framesPerChunk) / float64(sampleRate) * float64(time.Second))
	for off := 0; off < len(pcm); off += chunkSize {
		select {
		case e := <-recvErr:
			if e != nil && e != io.EOF {
				return e
			}
		default:
		}
		end := off + chunkSize
		if end > len(pcm) {
			end = len(pcm)
		}
		if e := stream.Send(pcm[off:end]); e != nil {
			return fmt.Errorf("send: %w", e)
		}
		time.Sleep(chunkDur)
	}
	stream.Close()
	select {
	case e := <-recvErr:
		if e != nil && e != io.EOF {
			return e
		}
	case <-time.After(2 * time.Second):
	}
	return nil
}

// readWAV returns the PCM data bytes and sample rate from a canonical PCM WAV.
func readWAV(path string) ([]byte, int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return nil, 0, fmt.Errorf("not a RIFF/WAVE file")
	}
	var sampleRate int
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
				sampleRate = int(binary.LittleEndian.Uint32(b[body+4 : body+8]))
			}
		case "data":
			pcm = b[body : body+size]
		}
		off = body + size
		if size%2 == 1 {
			off++ // chunks are word-aligned
		}
	}
	if sampleRate == 0 || pcm == nil {
		return nil, 0, fmt.Errorf("missing fmt/data chunk")
	}
	return pcm, sampleRate, nil
}

func shorten(s string) string {
	if i := len(s); i > 140 {
		return s[:137] + "..."
	}
	return s
}
