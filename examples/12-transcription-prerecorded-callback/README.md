# 12 — Transcribe prerecorded audio with callback

Set `Callback` on `PreRecordedTranscriptionOptions` to make the request async. `FromURL` returns immediately with a `request_id`; Deepgram POSTs the actual transcription to your URL when processing completes.

You need a publicly reachable HTTPS webhook to receive the result. Replace the placeholder URL in `main.go` before running. For local dev, expose a local server via ngrok / cloudflared / tailscale funnel.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
# Edit main.go to set Callback to your real webhook URL
go run ./examples/12-transcription-prerecorded-callback
```

## Equivalent Python

[`deepgram-python-sdk/examples/12-transcription-prerecorded-callback.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/12-transcription-prerecorded-callback.py)
