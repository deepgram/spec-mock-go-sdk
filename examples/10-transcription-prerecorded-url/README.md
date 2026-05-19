# 10 — Transcribe prerecorded audio from URL

`client.FromURL(ctx, url, opts)` POSTs `{"url": "..."}` to `/v1/listen` with the supplied options as query parameters. Returns the full parsed response.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/10-transcription-prerecorded-url
```

## See also

- [`15-transcription-advanced-options`](../15-transcription-advanced-options) — same call with diarization, smart format, etc.
- [`12-transcription-prerecorded-callback`](../12-transcription-prerecorded-callback) — async variant.
