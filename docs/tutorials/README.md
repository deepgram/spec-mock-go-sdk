# Tutorials

Step-by-step walk-throughs for the spec-driven Go SDK.

- [Transcribe your first file](./transcribe-your-first-file.md) — Listen REST, ~10 minutes, no setup beyond an API key.
- [Build a voice agent](./build-a-voice-agent.md) — Listen WebSocket streaming, ~20 minutes, demonstrates the production-resilience knobs.

Looking for something more granular than a tutorial? Every public symbol on the SDK has a runnable godoc example. See:

- [`pkg/client/listen/v1/rest/example_test.go`](../../pkg/client/listen/v1/rest/example_test.go) — 8 REST examples.
- [`pkg/client/listen/v1/websocket/example_test.go`](../../pkg/client/listen/v1/websocket/example_test.go) — 12 WS examples.

Migrating from the legacy `deepgram-go-sdk`? See [`../migration-from-deepgram-go-sdk.md`](../migration-from-deepgram-go-sdk.md).
