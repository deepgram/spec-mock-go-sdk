# Tutorial: Build a voice agent

End-to-end walk-through of the listen WebSocket surface. Streams microphone audio into Deepgram and prints live transcripts, including the production-resilience knobs you'd want for a real voice agent.

By the end you will:

- Have a Go program that streams audio and consumes events from `/v1/listen` over WebSocket.
- Understand the event-loop pattern (`Recv()` + type switch).
- Know how to apply backpressure, frame-size limits, and per-call timeouts.

## Prerequisites

- Go 1.24 or later.
- A Deepgram API key.
- An audio source. For this tutorial we generate synthetic 16-bit / 16 kHz PCM in code so the example is self-contained; in a real app you'd read from a microphone (PortAudio, `sox`, etc.) or a file.

## 1. Set up the project

```bash
mkdir voice-agent && cd voice-agent
go mod init example.com/voice-agent
go get github.com/deepgram/spec-mock-go-sdk@latest
export DEEPGRAM_API_KEY=your-key-here
```

## 2. Write the program

```go
package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func main() {
	client := wsv1.NewWithDefaults().WithConfig(&wsv1.Config{
		MaxFrameSizeBytes: 64 * 1024,
		SendTimeout:       5 * time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Connect(ctx, &wsv1.LiveTranscriptionOptions{
		Model:          "nova-3",
		Encoding:       "linear16",
		SampleRate:     16000,
		Channels:       1,
		InterimResults: true,
		SmartFormat:    true,
		Punctuate:      true,
	})
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer stream.Close()

	go pumpAudio(ctx, stream)

	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("stream ended")
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		switch e := event.(type) {
		case *wsv1.ResultsEvent:
			if len(e.Channel.Alternatives) > 0 {
				alt := e.Channel.Alternatives[0]
				if e.IsFinal {
					fmt.Println("FINAL:", alt.Transcript)
				} else if alt.Transcript != "" {
					fmt.Print("\rinterim: ", alt.Transcript)
				}
			}
		case *wsv1.SpeechStartedEvent:
			fmt.Println("[speech started]")
		case *wsv1.UtteranceEndEvent:
			fmt.Println("[utterance end]")
		case *wsv1.MetadataEvent:
			fmt.Println("[stream end, metadata received]")
			os.Exit(0)
		case *wsv1.ErrorEvent:
			log.Fatalf("stream error: %s", e.Description)
		case *wsv1.ReconnectEvent:
			fmt.Printf("[reconnected after %d attempts, %v backoff]\n",
				e.Attempts, e.BackoffApplied)
		}
	}
}

func pumpAudio(ctx context.Context, stream *wsv1.Stream) {
	// Simulate 5 seconds of silence followed by closing the stream.
	chunk := make([]byte, 3200) // 100ms of 16-bit / 16 kHz audio
	for i := 0; i < 50; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := stream.SendAudio(chunk); err != nil {
			log.Printf("send failed: %v", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := stream.CloseStream(); err != nil {
		log.Printf("close stream failed: %v", err)
	}
}

// generateSineWave fills buf with N samples of a sine wave at the
// supplied frequency. Useful for end-to-end smoke tests that don't
// have a real microphone.
func generateSineWave(buf []byte, freqHz, sampleRate int) {
	var w bytes.Buffer
	samples := len(buf) / 2
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		v := int16(math.Sin(2*math.Pi*float64(freqHz)*t) * 0.5 * math.MaxInt16)
		binary.Write(&w, binary.LittleEndian, v)
	}
	copy(buf, w.Bytes())
}
```

Save as `main.go`.

## 3. Run it

```bash
go run main.go
```

You will see `[speech started]` if the audio source is non-silent, `interim:` lines as Deepgram is forming results, `FINAL:` lines when a result is finalised, and finally `[stream end, metadata received]` when the session closes.

## What's happening

### The connect call

`client.Connect(ctx, opts)` dials the WebSocket and authenticates via the credentials on the client. The options become query parameters; the WebSocket upgrade carries them as part of the URL.

### The send pump

`pumpAudio` runs in its own goroutine, writing audio chunks via `stream.SendAudio`. The underlying transport serialises writes — multiple goroutines can call SendAudio concurrently without blocking each other beyond the write mutex.

### The receive loop

`stream.Recv()` blocks until the next server-emitted event. The returned `Event` is a tagged union; switch on its concrete type to handle each kind:

| Event variant | When it fires |
|---|---|
| `*ResultsEvent` | A transcript result is ready. Has `IsFinal` to distinguish interim updates from finalised text. |
| `*SpeechStartedEvent` | The server detected the beginning of an utterance. |
| `*UtteranceEndEvent` | The server detected the end of an utterance. |
| `*MetadataEvent` | The session is ending. Final summary metadata. |
| `*ErrorEvent` | The server reported an error and may close the connection. |
| `*SyncEvent` | The server echoed back a sync ID you sent via `stream.Sync(id)`. |
| `*ReconnectEvent` | The facade reconnected after a transport failure (only if `Config.ReconnectPolicy` is set). |

### The resilience knobs

`Config.MaxFrameSizeBytes` rejects oversized chunks before they reach the network. `Config.SendTimeout` bounds how long a single send will wait. Both are facade-only knobs that protect the SDK and your goroutines from getting wedged — they don't change what's on the wire.

## Common pitfalls

### Encoding mismatch

If your audio is not 16-bit linear PCM at 16 kHz, your transcript will be garbled. Set `Encoding`, `SampleRate`, and `Channels` to match your actual audio source. The defaults are correct only because we're sending silence in this tutorial.

### Forgetting to close the stream

A WebSocket session left dangling consumes resources on both sides. Always `defer stream.Close()`. To end the session gracefully and receive the final `MetadataEvent`, call `stream.CloseStream()` first; then continue receiving until you see `MetadataEvent` or `io.EOF`.

### Mixing Send and Recv goroutines

`stream.SendAudio` is goroutine-safe. `stream.Recv` is NOT — only one goroutine should call `Recv` at a time. The typical pattern (used above) is one send pump goroutine + the main loop driving Recv.

### Keep-alives in production

Idle connections are subject to driver / proxy timeouts. If you have stretches without audio (e.g. push-to-talk), send `stream.KeepAlive()` every 5-10 seconds. The wire shape is a small text frame and does not count toward billable audio.

## Where to go from here

- Add a real microphone source. PortAudio bindings or shelling to `sox`/`ffmpeg` are common.
- Set `Config.ReconnectPolicy` for transient-failure resilience.
- Pipe `*ResultsEvent.Channel.Alternatives[0].Transcript` into a downstream LLM for agent behaviour.
- Read [`pkg/client/listen/v1/websocket/example_test.go`](../../pkg/client/listen/v1/websocket/example_test.go) — every public symbol on the WS surface has a runnable godoc example.
