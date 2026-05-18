# Breaking changes

## TranscribeInput.Alternatives, Channels, SampleRate, FillerWords removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Members removed:** `Alternatives *int32`, `Channels *int32`, `SampleRate *int32`, `FillerWords *bool`
- **Customer-visible impact:** The facade options struct
  `PreRecordedTranscriptionOptions` still exposes `Alternatives`,
  `Channels`, `SampleRate`, and `FillerWords` so existing customer
  code keeps compiling. However, setting these fields is now a no-op:
  they are dropped in `optionsToTranscribeInput` and never reach the
  wire. Customers relying on these parameters reaching the Deepgram
  API will observe a behavioural change (server defaults apply).
- **Facade behaviour:** Wiring blocks for these four fields have been
  removed from `optionsToTranscribeInput`. `TestDropped_<Field>` tests
  in `wire_test.go` lock in the absorption to prevent silent re-wiring
  in a future regen.

## TranscribeInput hygiene removals (Uttseg, Dates, Name, EntityPrompt, SummarizeLength, Threshold, EmulateStreaming, Tier, UttSplitInterruptions, NumbersSpaces, KeywordBoost, Endpointing, Performance, Yelling, VadTurnon, DatasetId, Numbers, UnifySpeakerId, VadEvents, VadSuppression, Context, ShowRedactedText, Chunker, Identify, DateFormat, MaxSpeakers, Times, Ner)

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeInput`
- **Customer-visible impact:** None of these fields were exposed on
  the facade `PreRecordedTranscriptionOptions` struct nor wired in
  `optionsToTranscribeInput` prior to this regen. They were dropped
  from `api/types` and the facade simply no longer references them.
  Customer code is unaffected.
- **Facade behaviour:** No-op. No facade changes required for these.

## StreamInput.FillerWords, NoDelay, DiarizeModel removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Customer-visible impact:** The websocket facade does not have a
  `optionsToStreamInput` converter today — streaming options are
  serialised by the legacy `version.GetLiveAPI` query-string builder,
  not via the transport/http Invoke primitive. The facade
  `LiveTranscriptionOptions` struct still exposes `FillerWords`,
  `NoDelay`, and `DiarizeModel` for source-compatibility; whether
  they reach the wire depends on the legacy URL builder, which is
  out of scope for this regen.
- **Facade behaviour:** No facade source change required. Documented
  here so reviewers know the upstream removal is acknowledged.
