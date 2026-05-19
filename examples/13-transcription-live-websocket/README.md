# 13 — Live transcription via WebSocket

Streams a pre-recorded audio file in chunks to simulate a live microphone feed. Demonstrates the canonical send-pump + Recv-loop pattern: one goroutine pushes audio via `SendAudio`, the main goroutine drives `Recv` and type-switches over the `Event` variants.

Real microphone integrations swap the file-reading loop for a PortAudio binding or shelling to `sox`/`ffmpeg`.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/13-transcription-live-websocket
```

Uses [`../fixtures/audio.wav`](../fixtures/audio.wav) (16-bit mono 44.1 kHz PCM).

## See also

- [`pkg/client/listen/v1/websocket/resilience.go`](../../pkg/client/listen/v1/websocket/resilience.go) for production-resilience knobs (frame-size cap, send timeout, reconnect policy).
