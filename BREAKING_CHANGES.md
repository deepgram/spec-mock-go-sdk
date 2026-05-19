# Breaking changes

## Spec @internal hygiene audit (regen 2026-05-18)

The upstream Smithy spec (deepgram/spec PR #10) annotated a large set of fields
on `StreamInput` and `TranscribeInput` as `@internal`, removing them from the
public API surface. These fields were either undocumented stem parameters or
deprecated audit cleanup. They are now absent from the generated `api/types`
shapes.

### Tier 2 — unavoidable

For each removed field below: customers reaching into the generated wire-type
shape (`spectypes.StreamInput.*` or `spectypes.TranscribeInput.*`) will get a
compile error on next SDK upgrade. The facade `PreRecordedTranscriptionOptions`
does not expose any of these stem-only fields directly, so REST customers using
only `pkg/client/listen/v1/rest` are unaffected at the source level.

Customers using `LiveTranscriptionOptions` for the WebSocket streaming flow
may need to remove references to `FillerWords`, `NoDelay`, or `DiarizeModel`
on the streaming wire shape — these have been removed from `StreamInput`.
The prerecorded `TranscribeInput` equivalents (`FillerWords`, `DiarizeModel`)
remain and are still wired through the REST facade.

#### Removed from `spectypes.StreamInput`

- `FillerWords *bool` — never publicly documented for the streaming endpoint.
  The prerecorded `TranscribeInput.FillerWords` is unaffected.
- `NoDelay *bool` — internal-only flag, never publicly documented.
- `DiarizeModel *string` — internal-only on streaming; the prerecorded
  `TranscribeInput.DiarizeModel` is unaffected.

#### Removed from `spectypes.TranscribeInput`

All of the following were `@internal` / never-publicly-documented stem
parameters. None were exposed on the facade `PreRecordedTranscriptionOptions`,
so REST customers are not affected. Customers reaching directly into
`api/types` to set them will need to remove those references.

- `Uttseg *bool`
- `Dates *bool`
- `Name *string`
- `EntityPrompt *string`
- `SampleRate *int32`
- `SummarizeLength SummaryLength`
- `Threshold *int32`
- `EmulateStreaming *bool`
- `Tier *string`
- `UttSplitInterruptions *float32`
- `Alternatives *int32`
- `NumbersSpaces *bool`
- `KeywordBoost KeywordBoost`
- `Endpointing *string`
- `Performance *bool`
- `Yelling *bool`
- `Channels *int32`
- `VadTurnon *int32`
- `DatasetId *string`
- `Numbers *bool`
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

### Migration

If your code referenced any of the above on `spectypes.StreamInput` or
`spectypes.TranscribeInput` directly, remove the assignment. These options
were never documented at developers.deepgram.com and were not surfaced on
the idiomatic facade structs in `pkg/client/listen/v1/rest` or
`pkg/client/listen/v1/websocket`.

The facade-level `PreRecordedTranscriptionOptions` is unchanged at the
source level; the safety-net wire tests in `wire_test.go` continue to pass.
The facade-exposed `FillerWords` and `DiarizeModel` fields still wire through
to the prerecorded `TranscribeInput` shape, which retains those fields.

## StreamInput.FillerWords removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `FillerWords`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `FillerWords`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## StreamInput.NoDelay removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `NoDelay`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `NoDelay`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
- **Authored by:** the spec-idiomatic post-LLM tier-2 cleanup pass. This entry describes the actual code change applied in this run; it is not a forecast of how a production SDK ought to absorb a tier-2 removal. A future regen on a production SDK that elects a source-compat shim should write its own entry describing whatever shim it left behind.

## StreamInput.DiarizeModel removed

- **Tier:** 2 (unavoidable-breaking)
- **Shape:** `spectypes.StreamInput` (generated, in `api/types`)
- **Member:** `DiarizeModel`
- **Reason:** Removed from the public wire surface by upstream Smithy spec hygiene (e.g. `@internal` tagging or audit cleanup). The field was not part of the public customer-facing API.
- **What this PR did to the SDK:**
  - The facade did not expose `DiarizeModel`; no facade struct change applied.
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
  - The facade did not expose `SampleRate`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
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
  - The facade did not expose `Tier`; no facade struct change applied.
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
  - The facade did not expose `Alternatives`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
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
  - The facade did not expose `Endpointing`; no facade struct change applied.
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
  - The facade did not expose `Channels`; no facade struct change applied.
  - No matching `TestWires_/TestDropped_` test function existed.
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
  - The facade did not expose `VadEvents`; no facade struct change applied.
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

