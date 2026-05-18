# Breaking changes

This file enumerates spec changes that the facade in `pkg/` could not
fully absorb without a customer-visible behaviour shift, even though
customer source code keeps compiling.

## TranscribeInput.Alternatives removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Alternatives *int32`
- **Customer-visible impact:** The facade field
  `PreRecordedTranscriptionOptions.Alternatives int` is preserved so
  existing call sites still compile. The value is no longer wired
  through to the wire request; setting it has no effect on the
  transcription server.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.Alternatives`. `TestDropped_Alternatives` in `wire_test.go` locks
  in the absorption.

## TranscribeInput.Channels removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Channels *int32`
- **Customer-visible impact:** The facade field
  `PreRecordedTranscriptionOptions.Channels int` is preserved so
  existing call sites still compile. The value is no longer wired
  through to the wire request; setting it has no effect on the
  transcription server.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.Channels`. `TestDropped_Channels` in `wire_test.go` locks in the
  absorption.

## TranscribeInput.SampleRate removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `SampleRate *int32`
- **Customer-visible impact:** The facade field
  `PreRecordedTranscriptionOptions.SampleRate int` is preserved so
  existing call sites still compile. The value is no longer wired
  through to the wire request; setting it has no effect on the
  transcription server.
- **Facade behaviour:** `optionsToTranscribeInput` no longer references
  `in.SampleRate`. `TestDropped_SampleRate` in `wire_test.go` locks in
  the absorption.

## Other removed fields (unwired prior to this regen)

The following fields were removed from `TranscribeInput` and
`StreamInput` in this regen. None of them were wired through the
facade prior to the removal, so there is no behaviour change beyond
the wire-shape contraction:

- `StreamInput`: `FillerWords`, `NoDelay`, `DiarizeModel`
- `TranscribeInput`: `Uttseg`, `Dates`, `Name`, `EntityPrompt`,
  `SummarizeLength`, `Threshold`, `EmulateStreaming`, `Tier`,
  `UttSplitInterruptions`, `NumbersSpaces`, `KeywordBoost`,
  `Endpointing`, `Performance`, `Yelling`, `VadTurnon`, `DatasetId`,
  `Numbers`, `UnifySpeakerId`, `VadEvents`, `VadSuppression`,
  `Context`, `ShowRedactedText`, `Chunker`, `Identify`, `DateFormat`,
  `MaxSpeakers`, `Times`, `Ner`

These had no facade-options struct field wired through
`optionsToTranscribeInput` and no converter wiring; the facade
behaviour is unchanged.
