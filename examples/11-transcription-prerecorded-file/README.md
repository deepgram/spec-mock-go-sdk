# Prerecorded transcription from a file

This example sends a local file to Deepgram with `client.FromFile`.

It references the shared fixture at `examples/fixtures/audio.wav` and
passes `audio/wav` as the content type.

Run it from the repository root:

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/11-transcription-prerecorded-file
```
