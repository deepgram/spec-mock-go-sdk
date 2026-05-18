# Breaking changes

## TranscribeInput field removals (absorbed-with-ceremony)

The following fields were removed from `spectypes.TranscribeInput` in this
regen. The facade `interfaces.PreRecordedTranscriptionOptions` struct retains
the corresponding fields for source-compatibility — customer code that sets
them continues to compile — but `optionsToTranscribeInput` no longer wires
them to the wire shape. Setting these fields on the facade is now a silent
no-op; values do not reach the API.

- **Shape:** `spectypes.TranscribeInput`
- **Members removed:** `Alternatives`, `Channels`, `SampleRate`, `Uttseg`,
  `Dates`, `Name`, `EntityPrompt`, `SummarizeLength`, `Threshold`,
  `EmulateStreaming`, `Tier`, `UttSplitInterruptions`, `NumbersSpaces`,
  `KeywordBoost`, `Endpointing`, `Performance`, `Yelling`, `VadTurnon`,
  `DatasetId`, `Numbers`, `UnifySpeakerId`, `VadEvents`, `VadSuppression`,
  `Context`, `ShowRedactedText`, `Chunker`, `Identify`, `DateFormat`,
  `MaxSpeakers`, `Times`, `Ner`
- **Tier:** 1 (absorbed-breaking)
- **Customer-visible impact:** Any customer code that set one of these
  fields on `PreRecordedTranscriptionOptions` will continue to compile,
  but the value will no longer be sent to the API. Customers relying on
  these knobs (e.g. `Alternatives`, `Channels`, `SampleRate`) to control
  request behaviour will see the server fall back to defaults.
- **Facade behaviour:** `optionsToTranscribeInput` drops the wiring blocks
  for the previously-wired fields (`Alternatives`, `Channels`,
  `SampleRate`). The other fields in this list were never wired by the
  facade, so their removal from the spec is invisible at the wire
  boundary. The dropped-field list in `convert.go`'s docstring is updated;
  `TestDropped_Alternatives`, `TestDropped_Channels`, and
  `TestDropped_SampleRate` in `wire_test.go` lock in the absorption.

## StreamInput field removals (absorbed-with-ceremony)

- **Shape:** `spectypes.StreamInput` (live/streaming wire shape)
- **Members removed:** `FillerWords`, `NoDelay`, `DiarizeModel`
- **Tier:** 1 (absorbed-breaking)
- **Customer-visible impact:** None at the listen REST facade. These
  fields belong to the streaming/WebSocket surface. The listen WebSocket
  facade in `client/listen/v1/websocket/` does not currently wire a
  client→server options struct through a generated `StreamInput`, so no
  facade converter required adjustment. Customer code that configures
  live transcription via `LiveTranscriptionOptions` keeps the same
  fields — they are simply no longer sent to the server.
- **Facade behaviour:** No wiring drop required at the listen REST layer;
  documented here for audit completeness.
