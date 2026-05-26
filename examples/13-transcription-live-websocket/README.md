# Live transcription over WebSocket

This example uses `live.New` with the default WebSocket transport.

It connects to live transcription, sends audio chunks in one goroutine,
receives events in the main goroutine, then gracefully closes the stream.

Replace `readAudioChunks` with your microphone, file, or media pipeline.

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/13-transcription-live-websocket
```
