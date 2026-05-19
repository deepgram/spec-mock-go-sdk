# Examples

Runnable Go programs for common scenarios with the Deepgram Go SDK.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/10-transcription-prerecorded-url
```

Substitute any example directory name. All examples build via `go build ./examples/...` from the repo root.

## Index

### Authentication

| Example | Topic |
|---|---|
| [`01-authentication-api-key`](./01-authentication-api-key) | API-key auth via `DEEPGRAM_API_KEY` or explicit `New(...)`. |
| [`02-authentication-access-token`](./02-authentication-access-token) | Access-token auth via `DEEPGRAM_ACCESS_TOKEN`. |

### Transcription

| Example | Topic |
|---|---|
| [`10-transcription-prerecorded-url`](./10-transcription-prerecorded-url) | REST `FromURL` — transcribe audio at a public HTTPS URL. |
| [`11-transcription-prerecorded-file`](./11-transcription-prerecorded-file) | REST `FromFile` — transcribe a local audio file. |
| [`12-transcription-prerecorded-callback`](./12-transcription-prerecorded-callback) | REST with `Callback` — Deepgram POSTs back to your webhook. |
| [`13-transcription-live-websocket`](./13-transcription-live-websocket) | WebSocket streaming — chunked audio + typed event loop. |
| [`15-transcription-advanced-options`](./15-transcription-advanced-options) | Combined REST options — smart format + punctuate + diarize + language. |

### Configuration

| Example | Topic |
|---|---|
| [`70-request-options`](./70-request-options) | Passing arbitrary query parameters via `Extra`. |
| [`71-error-handling`](./71-error-handling) | Typed errors via `errors.As`, context cancellation, WebSocket error events. |

## Fixtures

`fixtures/audio.wav` — short PCM clip (16-bit mono 44.1 kHz). Used by `11` and `13`.
