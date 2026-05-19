# Breaking changes

## TranscribeInput / StreamInput: large @internal hygiene audit (spec PR #8)

- **Tier:** 2 (unavoidable)
- **Shapes:** `spectypes.TranscribeInput`, `spectypes.StreamInput`
- **Members removed from spec:** Uttseg, Dates, Name, EntityPrompt,
  SummarizeLength, Threshold, EmulateStreaming, Tier,
  UttSplitInterruptions, NumbersSpaces, KeywordBoost, Endpointing,
  Performance, Yelling, VadTurnon, DatasetId, Numbers, UnifySpeakerId,
  VadEvents, VadSuppression, Context, ShowRedactedText, Chunker,
  Identify, DateFormat, MaxSpeakers, Times, Ner (TranscribeInput);
  FillerWords, NoDelay, DiarizeModel (StreamInput).
- **Customer-visible impact:** Customers who set any of these fields on
  `interfaces.PreRecordedTranscriptionOptions` will see compile errors
  on this SDK upgrade. None of these fields were publicly documented
  on developers.deepgram.com; they were undocumented stem parameters.
- **Migration required:** Remove references to these fields. They had
  no observable effect for customers using the documented API surface.
- **Facade behaviour:** Facade option fields for the removed wire
  fields are dropped per the spec-hygiene policy (tier 2: surface the
  break to customers rather than silently no-op the setter).

## TranscribeInput.Alternatives / Channels / SampleRate

- **Tier:** 2 (unavoidable)
- **Shape:** `spectypes.TranscribeInput`
- **Customer-visible impact:** Facade options-struct fields kept for
  source-compat with existing customer code, but the values are
  silently dropped on the wire. Locked in by `TestDropped_*` cases in
  `wire_test.go`.
- **Migration required:** Stop setting these fields.
- **Facade behaviour:** Field retained on facade; converter wiring
  block dropped; documented in `convert.go`.
