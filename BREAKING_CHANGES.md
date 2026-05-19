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

## StreamInput.FillerWords removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `FillerWords`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `FillerWords` from its `*Options` struct.
  - Removed orphaned `TestWires_FillerWords` / `TestDropped_FillerWords` test function(s) that referenced the now-removed member.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## StreamInput.NoDelay removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `NoDelay`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `NoDelay` from its `*Options` struct.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## StreamInput.DiarizeModel removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `DiarizeModel`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `DiarizeModel` from its `*Options` struct.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Uttseg removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Uttseg`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Uttseg`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Dates removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Dates`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Dates`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Name removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Name`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Name`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.EntityPrompt removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `EntityPrompt`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `EntityPrompt`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.SampleRate removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `SampleRate`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `SampleRate` from its `*Options` struct.
  - Removed orphaned `TestWires_SampleRate` / `TestDropped_SampleRate` test function(s) that referenced the now-removed member.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.SummarizeLength removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `SummarizeLength`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `SummarizeLength`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Threshold removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Threshold`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Threshold`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.EmulateStreaming removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `EmulateStreaming`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `EmulateStreaming`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Tier removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Tier`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `Tier` from its `*Options` struct.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.UttSplitInterruptions removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `UttSplitInterruptions`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `UttSplitInterruptions`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Alternatives removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Alternatives`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `Alternatives` from its `*Options` struct.
  - Removed orphaned `TestWires_Alternatives` / `TestDropped_Alternatives` test function(s) that referenced the now-removed member.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.NumbersSpaces removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `NumbersSpaces`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `NumbersSpaces`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.KeywordBoost removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `KeywordBoost`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `KeywordBoost`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Endpointing removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Endpointing`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `Endpointing` from its `*Options` struct.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Performance removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Performance`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Performance`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Yelling removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Yelling`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Yelling`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Channels removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Channels`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `Channels` from its `*Options` struct.
  - Removed orphaned `TestWires_Channels` / `TestDropped_Channels` test function(s) that referenced the now-removed member.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.VadTurnon removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `VadTurnon`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `VadTurnon`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.DatasetId removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `DatasetId`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `DatasetId`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Numbers removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Numbers`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Numbers`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.UnifySpeakerId removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `UnifySpeakerId`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `UnifySpeakerId`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.VadEvents removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `VadEvents`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - Removed the facade member referencing `VadEvents` from its `*Options` struct.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.VadSuppression removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `VadSuppression`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `VadSuppression`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Context removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Context`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Context`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.ShowRedactedText removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `ShowRedactedText`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `ShowRedactedText`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Chunker removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Chunker`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Chunker`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Identify removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Identify`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Identify`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.DateFormat removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `DateFormat`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `DateFormat`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.MaxSpeakers removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `MaxSpeakers`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `MaxSpeakers`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Times removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Times`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Times`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## TranscribeInput.Ner removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.TranscribeInput` (generated, in `api/types`)
- **Member:** `Ner`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `Ner`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

