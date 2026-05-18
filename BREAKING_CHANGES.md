# Breaking changes

This file tracks breaking changes absorbed by the facade in response to upstream Smithy spec evolution. Entries are append-only; each describes the shape, member, customer-visible impact, and facade behaviour.

## TranscribeInput field removals (spec audit, multiple fields)

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Members removed:** `Alternatives`, `Channels`, `SampleRate`, `Uttseg`, `Dates`, `Name`, `EntityPrompt`, `SummarizeLength`, `Threshold`, `EmulateStreaming`, `Tier`, `UttSplitInterruptions`, `NumbersSpaces`, `KeywordBoost`, `Endpointing`, `Performance`, `Yelling`, `VadTurnon`, `DatasetId`, `Numbers`, `UnifySpeakerId`, `VadEvents`, `VadSuppression`, `Context`, `ShowRedactedText`, `Chunker`, `Identify`, `DateFormat`, `MaxSpeakers`, `Times`, `Ner`
- **Customer-visible impact:** Customer-facing `PreRecordedTranscriptionOptions` fields with these names continue to exist and accept values, but those values no longer reach the wire. For previously-wired fields (`Alternatives`, `Channels`, `SampleRate`), customers setting them will silently get default server behaviour instead of their chosen value. For the rest, behaviour is unchanged — the facade never wired them.
- **Facade behaviour:** `optionsToTranscribeInput` drops the wiring blocks for `Alternatives`, `Channels`, `SampleRate`. New `TestDropped_*` entries lock in the absorption so future regens can't silently re-wire fields the spec no longer models.

## StreamInput field removals

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Members removed:** `FillerWords`, `NoDelay`, `DiarizeModel`
- **Customer-visible impact:** Customer-facing `LiveTranscriptionOptions` continues to expose these fields. The streaming facade does not currently maintain a `streamingOptionsTo*` converter that wired these fields to the spec, so removal is a no-op on the active code path. Customers depending on these knobs reaching the wire (via custom-parameter contexts, etc.) should consult upstream docs for replacement parameters.
- **Facade behaviour:** No converter changes required; the facade options struct retains the fields for source-compat.
