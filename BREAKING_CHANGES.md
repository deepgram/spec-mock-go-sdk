# Breaking changes

## TranscribeInput.Alternatives removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Alternatives *int32`
- **Customer-visible impact:** Customer-facing `PreRecordedTranscriptionOptions.Alternatives int` field is preserved on the facade for source compatibility. Setting it has no effect on the wire request — the value is dropped on the floor.
- **Facade behaviour:** `optionsToTranscribeInput` no longer emits a wiring block for `Alternatives`. A `TestDropped_Alternatives` test in `wire_test.go` locks in the absorption.

## TranscribeInput.Channels removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Channels *int32`
- **Customer-visible impact:** Customer-facing `PreRecordedTranscriptionOptions.Channels int` is preserved; values set by callers are now silently dropped at the converter boundary.
- **Facade behaviour:** `optionsToTranscribeInput` no longer emits a wiring block for `Channels`. A `TestDropped_Channels` test locks in the absorption.

## TranscribeInput.SampleRate removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `SampleRate *int32`
- **Customer-visible impact:** Customer-facing `PreRecordedTranscriptionOptions.SampleRate int` is preserved; values set by callers are now silently dropped at the converter boundary.
- **Facade behaviour:** `optionsToTranscribeInput` no longer emits a wiring block for `SampleRate`. A `TestDropped_SampleRate` test locks in the absorption.

## Other @internal hygiene removals

The spec PR also removed a large number of additional `TranscribeInput` fields
(Uttseg, Dates, Name, EntityPrompt, SummarizeLength, Threshold,
EmulateStreaming, Tier, UttSplitInterruptions, NumbersSpaces, KeywordBoost,
Endpointing, Performance, Yelling, VadTurnon, DatasetId, Numbers,
UnifySpeakerId, VadEvents, VadSuppression, Context, ShowRedactedText, Chunker,
Identify, DateFormat, MaxSpeakers, Times, Ner) and `StreamInput` fields
(FillerWords, NoDelay, DiarizeModel) as part of an @internal audit. The facade
never wired any of these into the converter, so removal is a true no-op:
customer code that referenced them on the facade options struct (if any) will
still compile and the values will continue to be ignored on the wire — exactly
as they were before this regen.
