# Live transcription with SageMaker bidirectional streaming

This example routes live transcription through SageMaker HTTP/2
bidirectional streaming.

Import `github.com/aws/aws-sdk-go-v2/service/sagemakerruntimehttp2`
for this transport. Do not use `sagemakerruntime`; that package is for
standard request/response SageMaker calls.

The SDK keeps the same `Connect`, `SendAudio`, `Recv`, and `CloseStream`
streaming shape as the default WebSocket transport.
