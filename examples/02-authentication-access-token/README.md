# Authentication with an access token

This example constructs a prerecorded transcription client with
`prerecorded.WithAccessToken`.

Set `DEEPGRAM_ACCESS_TOKEN` before running it:

```bash
export DEEPGRAM_ACCESS_TOKEN="your_access_token"
go run ./examples/02-authentication-access-token
```

The client sends `Authorization: Bearer <token>` on the HTTP request.
