# 13 — Live transcription via WebSocket

Streams a pre-recorded audio file in chunks to simulate a live microphone feed. Demonstrates the canonical send-pump + Recv-loop pattern: one goroutine pushes audio via `SendAudio`, the main goroutine drives `Recv` and type-switches over the `Event` variants.

Real microphone integrations swap the file-reading loop for a PortAudio binding or shelling to `sox`/`ffmpeg`.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/13-transcription-live-websocket
```

Uses [`../fixtures/audio.wav`](../fixtures/audio.wav) (16-bit mono 44.1 kHz PCM).

## Equivalent Python

[`deepgram-python-sdk/examples/13-transcription-live-websocket.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/13-transcription-live-websocket.py)

| Python | Go |
|---|---|
| Event handlers via `connection.on(EventType.MESSAGE, ...)` | Type switch on `event := stream.Recv()` |
| `threading.Thread(target=connection.start_listening, ...)` | `go func() { ... stream.SendAudio(...) }()` |
| `connection.send_media(chunk)` | `stream.SendAudio(chunk)` |

## See also

- [`docs/tutorials/build-a-voice-agent.md`](../../docs/tutorials/build-a-voice-agent.md) for the full WebSocket walk-through.
- [`pkg/client/listen/v1/websocket/resilience.go`](../../pkg/client/listen/v1/websocket/resilience.go) for production-resilience knobs (frame-size cap, send timeout, reconnect policy).
