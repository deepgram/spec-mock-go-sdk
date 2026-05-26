# Prerecorded transcription from a URL

This example sends a hosted audio URL to Deepgram with `client.FromURL`.

It demonstrates the common prerecorded path: create a client, pass a
remote audio URL, and set typed transcription options.

Run it with:

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/10-transcription-prerecorded-url
```
