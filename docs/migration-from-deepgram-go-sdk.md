# Migrating from `github.com/deepgram/deepgram-go-sdk`

This is the spec-driven Go SDK. It lives at a new module path
(`github.com/deepgram/sdk-go` once published — see [ADR-0001](./adrs/0001-module-path.md))
and runs side-by-side with the legacy `github.com/deepgram/deepgram-go-sdk`.

**You do not need to migrate.** Your existing code on `deepgram-go-sdk` continues to work. The legacy SDK stays maintained on its current import path for the foreseeable future. This document is for customers who want to opt-in to the new SDK after weighing what they gain and what they pay.

## Should you migrate?

### Migrate if:

- You want the spec-driven generation pipeline's benefits: faster lockstep with the public Deepgram surface, deprecations propagated automatically, removed-from-spec parameters surfaced at compile time rather than runtime.
- You want a smaller dependency footprint. The new SDK has 3 direct deps (websocket, AWS sagemaker, smithy-go) vs the legacy SDK's broader surface.
- You're starting a new Go project against Deepgram and have no existing import-path commitment.
- You're rebuilding your integration layer anyway and the migration cost is amortised.

### Stay on `deepgram-go-sdk` if:

- Your integration is stable and the migration cost outweighs the benefits.
- You depend on legacy-SDK behaviour that has not yet shipped in this SDK (see "Coverage gaps" below).
- You're not Go-idiomatic and prefer the legacy SDK's existing patterns.
- You're risk-averse and want to wait for v1.0.0.

## Coverage gaps

The spec-driven SDK is currently in alpha. It ships with:

- **Listen REST** (`POST /v1/listen`) — full @httpQuery surface for `pkg/client/listen/v1/rest`.
- **Listen WebSocket** (`/v1/listen` over WS) — full streaming surface for `pkg/client/listen/v1/websocket`, including the resilience knobs documented in `resilience.go`.

It does NOT yet ship:

- Speak (TTS) — coming.
- Agent — coming.
- Manage / projects / keys / billing — coming.
- Read / analyze / pre-recorded analyze — coming.
- Self-hosted licensing helpers.

Track each on the [`deepgram/spec-mock-go-sdk` issues board](https://github.com/deepgram/spec-mock-go-sdk/issues).

## What's different

### Import paths

```go
// Legacy
import (
    "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
    "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"
)

// Spec-driven (this SDK)
import (
    restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
    wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)
```

The package layout intentionally mirrors the spec's per-product / per-version / per-transport tree, so future products land at predictable paths.

### Options structs

The new SDK uses one Options struct per transport per operation, named after the customer-visible concept:

| Legacy | Spec-driven |
|---|---|
| Mixed transport options on a single client struct | `restv1.PreRecordedTranscriptionOptions` for REST batch |
| | `wsv1.LiveTranscriptionOptions` for WebSocket streaming |
| | Plus `wsv1.Config` for facade-level resilience knobs (frame-size limits, send timeouts, reconnect policy) |

Fields use idiomatic Go zero-value semantics: leave a field at its zero value to omit it from the request. No optional-pointer wrapping.

### Constructor patterns

```go
// Legacy (illustrative; check legacy README for exact)
client := listen.New(deepgram.ClientOptions{ApiKey: "..."})

// Spec-driven
client := restv1.New("api-key", "")
// or
client := restv1.NewWithDefaults()
// or
client := restv1.NewWithDefaults().WithBaseURL("https://staging.api.deepgram.com")
```

`NewWithDefaults()` reads `DEEPGRAM_API_KEY` / `DEEPGRAM_ACCESS_TOKEN` from the environment. `WithBaseURL` / `WithConfig` / `WithHTTPClient` are immutable-builder methods returning copies.

### Removed / renamed parameters

The spec-driven SDK strips `@internal`-tagged parameters that were never publicly documented but appeared on the legacy surface. If you were sending one of these on `?...`, the new options struct will not expose the field; the parameter is no longer accepted by your code at compile time.

Notable removals at first ship:
- (Tracked in `BREAKING_CHANGES.md` once the file accumulates real entries; the rehearsal entry for `Dictation` is illustrative-only and reverts before each release.)

### `@deprecated` markers

Deprecated parameters render with `// Deprecated:` godoc blocks. Your IDE surfaces these inline; `go vet -staticcheck` flags them at build time. The legacy SDK did not consistently propagate deprecations from the public surface; the new SDK does.

### WebSocket session lifecycle

```go
// Spec-driven WebSocket pattern
stream, err := wsClient.Connect(ctx, &wsv1.LiveTranscriptionOptions{
    Model:          "nova-3",
    Encoding:       "linear16",
    SampleRate:     16000,
    InterimResults: true,
})
if err != nil { return err }
defer stream.Close()

// Send audio
for {
    n, err := mic.Read(buf)
    if err != nil { break }
    if err := stream.SendAudio(buf[:n]); err != nil { return err }
}

// Receive events
for {
    event, err := stream.Recv()
    if err == io.EOF { return nil }
    if err != nil { return err }
    switch e := event.(type) {
    case *wsv1.ResultsEvent:
        // ...
    case *wsv1.MetadataEvent:
        return nil
    case *wsv1.ErrorEvent:
        return fmt.Errorf("stream error: %s", e.Description)
    }
}
```

The receive loop yields typed `Event` variants (`*ResultsEvent`, `*MetadataEvent`, `*SpeechStartedEvent`, `*UtteranceEndEvent`, `*ErrorEvent`, `*SyncEvent`) plus the facade-only `*ReconnectEvent`. Legacy callback-style handlers are not provided; if you need them, wrap `stream.Recv()` in a small dispatcher.

## What's the same

- **Credentials**. `DEEPGRAM_API_KEY` and `DEEPGRAM_ACCESS_TOKEN` env vars work identically. `Authorization: Token <api-key>` and `Authorization: Bearer <access-token>` go on the wire the same way.
- **Wire surface**. The HTTP routes, query parameters, response shapes, and WebSocket frame formats are the same as what the Deepgram service has always accepted. The SDK changes how YOU talk to the SDK, not how the SDK talks to the service.
- **Self-hosted endpoints**. Use `WithBaseURL("https://your-host:port")` (or `wss://...` for the WebSocket client).
- **`go get`-able**. Once the new module path publishes, `go get github.com/deepgram/sdk-go@latest` works the same way as the legacy SDK.

## Common migration recipe

1. **Add the new SDK as a dep alongside the legacy one.**

   ```bash
   go get github.com/deepgram/spec-mock-go-sdk@latest
   ```

   (or `github.com/deepgram/sdk-go@latest` once renamed)

2. **Migrate one product at a time.** Start with whichever has the lightest integration surface (typically the pre-recorded REST path). Leave the rest on the legacy SDK.

3. **Move credential handling to the new SDK's constructors.** `restv1.NewWithDefaults()` is the simplest path.

4. **Adapt option-setting code to the per-operation Options structs.** This is the bulk of the migration work. Most fields have 1:1 names; check `pkg/client/listen/v1/rest/options.go` for the exhaustive list.

5. **Handle the receive loop change for WebSocket.** Move from callback handlers (if you used them) to the `Recv()` + type-switch pattern.

6. **Run your test suite.** The new SDK ships wire tests for every facade field; if your test suite has end-to-end coverage, it should pass after the option-struct rewrites.

7. **Remove the legacy SDK dep when no calls remain.** `go mod tidy` will surface remaining references.

## Codemod

A `gofmt`-style migration tool (`deepgram/sdk-go-migrate`) is on the work plan but not shipped. When it lands, it will:

- Rewrite import paths from `deepgram/deepgram-go-sdk` to the new module path.
- Translate the most common option-setting patterns.
- Flag remaining manual-migration items via `// MIGRATE:` comments.

Until then, the migration is manual but bounded. Most integrations touch <5 files.

## Questions / problems

Open an issue on [`deepgram/spec-mock-go-sdk`](https://github.com/deepgram/spec-mock-go-sdk/issues) with the `migration` label. Include:

- Legacy code snippet.
- What you expected the new SDK to do.
- What actually happened (error, behaviour difference, surface gap).

For coverage gaps (a product or option not yet in the spec-driven SDK), drop the issue and stay on the legacy SDK for that path. We're prioritising spec coverage by customer demand.
