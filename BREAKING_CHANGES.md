# Breaking changes

## Spec @internal hygiene audit (regen 2026-05-18)

The upstream Smithy spec (deepgram/spec PR #10) annotated a large set of fields
on `StreamInput` and `TranscribeInput` as `@internal`, removing them from the
public API surface. These fields were either undocumented stem parameters or
deprecated audit cleanup. They are now absent from the generated `api/types`
shapes.

### Tier 2 — unavoidable

For each removed field below: customers reaching into the generated wire-type
shape (`spectypes.StreamInput.*` or `spectypes.TranscribeInput.*`) will get a
compile error on next SDK upgrade. The facade `PreRecordedTranscriptionOptions`
does not expose any of these fields directly, so REST customers using only
`pkg/client/listen/v1/rest` are unaffected at the source level.

Customers using `LiveTranscriptionOptions` for the WebSocket streaming flow
may need to remove references to `FillerWords`, `NoDelay`, or `DiarizeModel`
if those fields were being set there — these have been removed from the
streaming wire shape.

#### Removed from `spectypes.StreamInput`

- `FillerWords *bool` — never publicly documented for the streaming endpoint;
  facade `LiveTranscriptionOptions.FillerWords` if present is silently dropped.
- `NoDelay *bool` — internal-only flag, never publicly documented.
- `DiarizeModel *string` — internal-only; legacy `Diarize`+`DiarizeVersion` remain.

#### Removed from `spectypes.TranscribeInput`

All of the following were `@internal` / never-publicly-documented stem
parameters. None were exposed on the facade `PreRecordedTranscriptionOptions`,
so REST customers are not affected. Customers reaching directly into
`api/types` to set them will need to remove those references.

- `Uttseg *bool`
- `Dates *bool`
- `Name *string`
- `EntityPrompt *string`
- `SampleRate *int32`
- `SummarizeLength SummaryLength`
- `Threshold *int32`
- `EmulateStreaming *bool`
- `Tier *string`
- `UttSplitInterruptions *float32`
- `Alternatives *int32`
- `NumbersSpaces *bool`
- `KeywordBoost KeywordBoost`
- `Endpointing *string`
- `Performance *bool`
- `Yelling *bool`
- `Channels *int32`
- `VadTurnon *int32`
- `DatasetId *string`
- `Numbers *bool`
- `UnifySpeakerId *bool`
- `VadEvents *bool`
- `VadSuppression *bool`
- `Context *bool`
- `ShowRedactedText *bool`
- `Chunker Chunker`
- `Identify *bool`
- `DateFormat *string`
- `MaxSpeakers *int32`
- `Times *bool`
- `Ner *bool`

### Migration

If your code referenced any of the above on `spectypes.StreamInput` or
`spectypes.TranscribeInput` directly, remove the assignment. These options
were never documented at developers.deepgram.com and were not surfaced on
the idiomatic facade structs in `pkg/client/listen/v1/rest` or
`pkg/client/listen/v1/websocket`.

The facade-level `PreRecordedTranscriptionOptions` is unchanged at the
source level; the safety-net wire tests in `wire_test.go` continue to pass.
