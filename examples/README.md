# Examples

Runnable examples for the spec-driven Go SDK. Organised by feature area mirroring the [`deepgram-python-sdk` examples convention](https://github.com/deepgram/deepgram-python-sdk/tree/main/examples) — each section starts at a multiple of 10.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/10-transcription-prerecorded-url
```

Substitute any example directory name. All examples build via `go build ./examples/...` from the repo root.

## Index

### 01-09 Authentication

| Example | Topic |
|---|---|
| [`01-authentication-api-key`](./01-authentication-api-key) | API-key auth via `DEEPGRAM_API_KEY` or explicit `New(...)`. |
| [`02-authentication-access-token`](./02-authentication-access-token) | Access-token auth via `DEEPGRAM_ACCESS_TOKEN`. (Token issuance is out of scope for this SDK — typically dx-id provisioned.) |

### 10-19 Transcription (Listen)

| Example | Topic |
|---|---|
| [`10-transcription-prerecorded-url`](./10-transcription-prerecorded-url) | REST `FromURL` — transcribe audio at a public HTTPS URL. |
| [`11-transcription-prerecorded-file`](./11-transcription-prerecorded-file) | REST `FromFile` — transcribe a local audio file. |
| [`12-transcription-prerecorded-callback`](./12-transcription-prerecorded-callback) | REST with `Callback` — Deepgram POSTs back to your webhook. |
| [`13-transcription-live-websocket`](./13-transcription-live-websocket) | WebSocket streaming — chunked audio + typed event loop. |
| [`15-transcription-advanced-options`](./15-transcription-advanced-options) | Combined REST options — smart format + punctuate + diarize + language. |

### 70-79 Configuration & Advanced

| Example | Topic |
|---|---|
| [`70-request-options`](./70-request-options) | Idiomatic typed options vs Python's `additional_query_parameters` escape hatch. |
| [`71-error-handling`](./71-error-handling) | Errors from REST + WS, `errors.Is`, context cancellation, `ErrFrameTooLarge` / `ErrSendTimeout`. |

## Not yet covered

| Python example | Why it's absent |
|---|---|
| `14-transcription-live-websocket-v2.py` | Listen V2 is not in the Smithy spec yet. |
| `20-29` Speak (TTS) | Speak product not in this prototype. |
| `27-transcription-live-sagemaker.py` | SageMaker transport scaffolded in `api/transport/sagemaker` but no facade. |
| `30 Voice Agent` | Agent product not in this prototype. |
| `40 Text Intelligence` | Read product not in this prototype. |
| `50-56 Management` | Management product not in this prototype. |
| `60 On-Premises` | Self-hosted credentials story not in this prototype. |

These return as the spec coverage expands.

## Fixtures

`fixtures/audio.wav` — short PCM clip (16-bit mono 44.1 kHz, mirrors the Python SDK fixture). Used by `11` and `13`.
