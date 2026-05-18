# Breaking changes

## TranscribeInput.Alternatives removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Alternatives *int32`
- **Customer-visible impact:** `PreRecordedTranscriptionOptions.Alternatives` is retained on the facade for source-compatibility, but values set on it no longer reach the wire. Requests will use the server default.
- **Facade behaviour:** `optionsToTranscribeInput` no longer wires `Alternatives`. `TestDropped_Alternatives` locks in the absorption.

## TranscribeInput.Channels removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput`
- **Member:** `Channels *int32`
- **Customer-visible impact:** `PreRecordedTranscriptionOptions.Channels` setter is now a no-op on the wire. The server infers channel count from the audio.
- **Facade behaviour:** wiring dropped from `optionsToTranscribeInput`; `TestDropped_Channels` enforces the drop.

## TranscribeInput.SampleRate removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput`
- **Member:** `SampleRate *int32`
- **Customer-visible impact:** `PreRecordedTranscriptionOptions.SampleRate` setter is now a no-op on the wire. Sample rate is inferred from the audio container.
- **Facade behaviour:** wiring dropped from `optionsToTranscribeInput`; `TestDropped_SampleRate` enforces the drop.

## TranscribeInput: bulk removal of unwired fields

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput`
- **Members:** `Uttseg`, `Dates`, `Name`, `EntityPrompt`, `SummarizeLength`, `Threshold`, `EmulateStreaming`, `Tier`, `UttSplitInterruptions`, `NumbersSpaces`, `KeywordBoost`, `Endpointing`, `Performance`, `Yelling`, `VadTurnon`, `DatasetId`, `Numbers`, `UnifySpeakerId`, `VadEvents`, `VadSuppression`, `Context`, `ShowRedactedText`, `Chunker`, `Identify`, `DateFormat`, `MaxSpeakers`, `Times`, `Ner`
- **Customer-visible impact:** None of these fields were wired by the facade prior to the regen. Customer code that set the corresponding `PreRecordedTranscriptionOptions` fields (where they exist) already had no effect on the wire. No additional regression.
- **Facade behaviour:** No wiring changes required; the converter never referenced these.

## StreamInput: FillerWords, NoDelay, DiarizeModel removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.StreamInput`
- **Members:** `FillerWords`, `NoDelay`, `DiarizeModel`
- **Customer-visible impact:** Streaming (WebSocket) options struct still exposes these fields for source-compat. They no longer reach the streaming URL query string when this regen applies; behaviour falls back to server defaults.
- **Facade behaviour:** Streaming URL builder absorbs the absence by skipping unknown fields. No converter code references `StreamInput` directly in the facade.
