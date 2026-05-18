# Breaking changes

This file is maintained by `spec-idiomatic` on every regen. Each entry
describes a single shape/member change classified by tier per
`.agents/skills/sdk-breaking-ceremony/SKILL.md`.

## TranscribeInput.Alternatives removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Alternatives *int32`
- **Customer-visible impact:** The facade option
  `PreRecordedTranscriptionOptions.Alternatives int` is preserved so
  customer code keeps compiling, but it no longer reaches the wire.
  Requests that previously set `Alternatives` to influence the number
  of ranked interpretations returned will be sent without that query
  parameter and the server default will apply.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.Alternatives`. `TestDropped_Alternatives` locks in the absorption.

## TranscribeInput.Channels removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Channels *int32`
- **Customer-visible impact:** The facade option
  `PreRecordedTranscriptionOptions.Channels int` is preserved but no
  longer wired. Customers setting `Channels` to declare interleaved
  channel count for raw audio uploads will silently rely on the server
  default.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.Channels`. `TestDropped_Channels` locks in the absorption.

## TranscribeInput.SampleRate removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `SampleRate *int32`
- **Customer-visible impact:** The facade option
  `PreRecordedTranscriptionOptions.SampleRate int` is preserved but no
  longer wired. Customers setting `SampleRate` to declare raw audio
  sample rate will silently rely on the server default.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.SampleRate`. `TestDropped_SampleRate` locks in the absorption.

## Other TranscribeInput / StreamInput fields removed in this regen

The spec also removed the following `TranscribeInput` fields, none of
which were wired through the facade prior to this regen, so the
absorption is a no-op on the customer surface (the converter had no
`if o.X` block for them):

- `Uttseg`, `Dates`, `Name`, `EntityPrompt`, `SummarizeLength`,
  `Threshold`, `EmulateStreaming`, `Tier`, `UttSplitInterruptions`,
  `NumbersSpaces`, `KeywordBoost`, `Endpointing`, `Performance`,
  `Yelling`, `VadTurnon`, `DatasetId`, `Numbers`, `UnifySpeakerId`,
  `VadEvents`, `VadSuppression`, `Context`, `ShowRedactedText`,
  `Chunker`, `Identify`, `DateFormat`, `MaxSpeakers`, `Times`, `Ner`

The streaming `StreamInput` fields `FillerWords`, `NoDelay`, and
`DiarizeModel` were also removed; the streaming facade does not
currently expose those knobs, so no customer-visible change.
