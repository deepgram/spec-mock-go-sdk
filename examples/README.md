# Examples

Standalone programs for the Deepgram Go SDK speech-to-text packages.

| Example | Description |
|---|---|
| [`01-authentication-api-key`](01-authentication-api-key) | Authenticate with a Deepgram API key. |
| [`02-authentication-access-token`](02-authentication-access-token) | Authenticate with a Bearer access token. |
| [`10-transcription-prerecorded-url`](10-transcription-prerecorded-url) | Transcribe audio from a hosted URL. |
| [`11-transcription-prerecorded-file`](11-transcription-prerecorded-file) | Transcribe a local audio file. |
| [`12-transcription-prerecorded-callback`](12-transcription-prerecorded-callback) | Request an async callback for prerecorded transcription. |
| [`13-transcription-live-websocket`](13-transcription-live-websocket) | Stream live audio over the default WebSocket transport. |
| [`15-transcription-advanced-options`](15-transcription-advanced-options) | Combine typed options with additional query parameters. |
| [`20-transcription-prerecorded-sagemaker`](20-transcription-prerecorded-sagemaker) | Use SageMaker for prerecorded transcription. |
| [`21-transcription-live-sagemaker-bidi`](21-transcription-live-sagemaker-bidi) | Use SageMaker HTTP/2 bidirectional streaming for live transcription. |
| [`70-request-options`](70-request-options) | Send parameters that are not yet typed in the SDK. |
| [`71-error-handling`](71-error-handling) | Handle transport and typed API errors. |

Build all examples with:

```bash
go build ./examples/...
```
