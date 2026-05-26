# Prerecorded transcription with a callback

This example sets `Callback` and `CallbackMethod` on
`PreRecordedTranscriptionOptions`.

Deepgram can return quickly and send the completed transcription to
your callback URL when processing finishes.

Run it with your own publicly reachable callback endpoint:

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/12-transcription-prerecorded-callback
```
