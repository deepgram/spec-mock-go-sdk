# Authentication with an API key

This example constructs a prerecorded transcription client with
`prerecorded.WithAPIKey`.

Set `DEEPGRAM_API_KEY` before running it:

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/01-authentication-api-key
```

The client sends `Authorization: Token <key>` on the HTTP request.
