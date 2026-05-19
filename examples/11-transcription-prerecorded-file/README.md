# 11 — Transcribe prerecorded audio from file

`client.FromFile(ctx, path, contentType, opts)` streams a local audio file as the request body to `/v1/listen`.

For in-memory bytes or non-file readers, swap to `client.FromStream(ctx, reader, contentType, opts)`.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/11-transcription-prerecorded-file
```

Uses [`../fixtures/audio.wav`](../fixtures/audio.wav) (16-bit mono 44.1 kHz PCM).

## Equivalent Python

[`deepgram-python-sdk/examples/11-transcription-prerecorded-file.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/11-transcription-prerecorded-file.py)

| Python | Go |
|---|---|
| `client.listen.v1.media.transcribe_file(request=audio_data, model=...)` (in-memory bytes) | `client.FromStream(ctx, reader, contentType, opts)` |
| Hand-roll `read_file_in_chunks(...)` generator | `client.FromFile(ctx, path, contentType, opts)` (streams from disk internally) |
