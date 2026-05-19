# 70 — Request options

## Spec-driven, not arbitrary-param

The Python SDK exposes a `request_options` dict with an `additional_query_parameters` slot for passing arbitrary query parameters Deepgram might accept that the SDK does not have a typed field for.

This Go SDK is **spec-driven**: every option is a typed field on `PreRecordedTranscriptionOptions` corresponding to a `@httpQuery` member on the Smithy spec's `TranscribeInput`. There is no arbitrary-param escape hatch by design — adding a new field requires:

1. Add the `@httpQuery` member to `deepgram/spec`.
2. Regen via `tools/codegen.sh` → `api/types.go` gets the new field.
3. `spec-idiomatic` adds the matching field to `PreRecordedTranscriptionOptions` (via the regen loop).

The trade-off: customer code can never desync from the documented Deepgram surface, but the SDK cannot expose undocumented experimental parameters without a spec roll-forward.

## What this example shows

The most common reason to reach for arbitrary params in the Python SDK is `DetectLanguage`, which IS a typed field here. The example demonstrates passing it just like any other option.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/70-request-options
```

## Equivalent Python

[`deepgram-python-sdk/examples/70-request-options.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/70-request-options.py)

| Python | Go |
|---|---|
| `request_options={"additional_query_parameters": {"detect_language": ["en", "es"]}}` | `DetectLanguage: []string{"en", "es"}` |
| `"timeout_in_seconds": 30` | Use `context.WithTimeout(...)` on the call's context. |
| `"max_retries": 3` | Wrap the call site in your own retry loop. The SDK does not retry. |
| `"additional_headers": {"X-Custom-Header": "..."}` | Pass a custom `*http.Client` via `client.WithHTTPClient(...)` whose `RoundTripper` adds the header. |
