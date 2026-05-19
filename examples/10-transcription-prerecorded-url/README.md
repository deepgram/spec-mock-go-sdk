# 10 — Transcribe prerecorded audio from URL

`client.FromURL(ctx, url, opts)` POSTs `{"url": "..."}` to `/v1/listen` with the supplied options as query parameters. Returns the full parsed response.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/10-transcription-prerecorded-url
```

## Equivalent Python

[`deepgram-python-sdk/examples/10-transcription-prerecorded-url.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/10-transcription-prerecorded-url.py)

| Python | Go |
|---|---|
| `client.listen.v1.media.transcribe_url(url=..., model=...)` | `client.FromURL(ctx, url, &restv1.PreRecordedTranscriptionOptions{Model: "..."})` |

## See also

- [`15-transcription-advanced-options`](../15-transcription-advanced-options) — same call with diarization, smart format, etc.
- [`12-transcription-prerecorded-callback`](../12-transcription-prerecorded-callback) — async variant.
