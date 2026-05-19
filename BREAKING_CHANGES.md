# Breaking changes

This file is auto-maintained by the `spec-idiomatic` runner. Each entry
describes a spec change that could not be silently absorbed by the
facade.

## TranscribeInput: large @internal hygiene audit (spec PR #8)

The upstream Smithy spec removed a long list of fields from
`TranscribeInput` as part of the `@internal` annotation alignment with
curated docs. These fields were either internal-only stem parameters,
undocumented on developers.deepgram.com, or deprecated. The runner
classified each as **tier 2 / unavoidable** and applied the spec-hygiene
removal policy: the corresponding facade option fields are dropped from
`PreRecordedTranscriptionOptions` rather than kept as silent no-ops, so
customer code that referenced these fields will fail to compile on this
SDK upgrade.

Customers using only documented options on
https://developers.deepgram.com/reference/pre-recorded should be
unaffected. Customers reaching into undocumented stem parameters will
see a compile break — that is the intended ceremony.

### TranscribeInput fields removed (unavoidable, facade field dropped)

- **Tier:** 2 (unavoidable)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Members removed from spec AND facade:**
  - `Uttseg *bool`
  - `Dates *bool`
  - `Name *string`
  - `EntityPrompt *string`
  - `SummarizeLength SummaryLength`
  - `Threshold *int32`
  - `EmulateStreaming *bool`
  - `Tier *string`
  - `UttSplitInterruptions *float32`
  - `NumbersSpaces *bool`
  - `KeywordBoost KeywordBoost`
  - `Endpointing *string`
  - `Performance *bool`
  - `Yelling *bool`
  - `VadTurnon *int32`
  - `DatasetId *string`
  - `Numbers *bool` (superseded by `Numerals`)
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
- **Customer-visible impact:** Any code that set these fields on
  `interfaces.PreRecordedTranscriptionOptions` (e.g.
  `&interfaces.PreRecordedTranscriptionOptions{Tier: "enhanced"}`)
  will fail to compile.
- **Migration required:** Remove references to these fields. They were
  never publicly documented and had no observable effect for customers
  using the documented API surface.
- **Facade behaviour:** Fields removed from
  `PreRecordedTranscriptionOptions`; no wiring blocks exist in
  `optionsToTranscribeInput`.

### StreamInput fields removed (unavoidable)

- **Tier:** 2 (unavoidable)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Members removed:**
  - `FillerWords *bool`
  - `NoDelay *bool`
  - `DiarizeModel *string`
- **Customer-visible impact:** None for the listen WS facade — these
  fields were never exposed on `LiveTranscriptionOptions` in the
  facade. The wire input shape no longer carries them; customers
  calling the live API without setting these are unaffected.
- **Facade behaviour:** No facade options struct exposed these as
  wired members. No changes to facade.

### TranscribeInput.Alternatives / Channels / SampleRate (unavoidable, facade field retained)

- **Tier:** 2 (unavoidable)
- **Shape:** `spectypes.TranscribeInput`
- **Members:** `Alternatives`, `Channels`, `SampleRate` (all removed
  from the generated wire shape)
- **Customer-visible impact:** Setting these fields on
  `PreRecordedTranscriptionOptions` compiles but the values are
  silently dropped on the wire. Locked in by `TestDropped_Alternatives`,
  `TestDropped_Channels`, `TestDropped_SampleRate` in `wire_test.go`.
- **Migration required:** Stop setting these fields; they are not part
  of the documented public API. The fields are retained on the facade
  options struct to preserve source-compat for existing callers; the
  values simply do not reach the wire.
- **Facade behaviour:** Facade options-struct field kept; the wiring
  block was dropped because the generated `TranscribeInput` no longer
  carries these fields. Documented in `convert.go`.
