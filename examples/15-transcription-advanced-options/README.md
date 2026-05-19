# 15 — Transcription with advanced options

Combines the option fields most production integrations want — smart formatting, punctuation, diarization, explicit language pin — in a single `FromURL` call.

The full inventory lives in [`pkg/client/listen/v1/rest/options.go`](../../pkg/client/listen/v1/rest/options.go). Each field corresponds to one `@httpQuery` member on the Smithy spec's `TranscribeInput`.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/15-transcription-advanced-options
```

## Equivalent Python

[`deepgram-python-sdk/examples/15-transcription-advanced-options.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/15-transcription-advanced-options.py)

| Python | Go |
|---|---|
| `smart_format=True` | `SmartFormat: true` |
| `punctuate=True` | `Punctuate: true` |
| `diarize=True` | `Diarize: true` |
| `language="en-US"` | `Language: "en-US"` |

The Go SDK uses CamelCase field names and Go zero-value semantics: leave a field unset to omit it from the request. No `Optional[bool]` wrapping needed.
