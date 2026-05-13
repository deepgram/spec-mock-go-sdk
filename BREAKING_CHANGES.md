# Breaking changes

## TranscribeOutput.Metadata removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeOutput` (generated, in `api/types`)
- **Member:** `Metadata *ResponseMetadata`
- **Customer-visible impact:** The customer-facing
  `PreRecordedResponse.Metadata *Metadata` field in
  `pkg/client/listen/v1/rest` is preserved in the facade so the Go API
  signature stays stable, but it will now always be `nil`. Any customer
  code that read fields off `resp.Metadata` (e.g. `resp.Metadata.RequestID`,
  `resp.Metadata.Duration`, `resp.Metadata.ModelInfo`,
  `resp.Metadata.Warnings`, `resp.Metadata.SummaryInfo`, etc.) will see
  zero values or nil-dereference unless it nil-checks `resp.Metadata`
  first. `resp.RequestID` (top-level on `PreRecordedResponse`) is
  unaffected and continues to be populated from
  `spectypes.TranscribeOutput.RequestId`.
- **Facade behaviour:** `convertTranscribeOutput` in
  `pkg/client/listen/v1/rest/convert.go` no longer references
  `in.Metadata`; it hardcodes `Metadata: nil` in the returned
  `*PreRecordedResponse`. The `convertResponseMetadata` helper has been
  removed since it has no remaining call sites.
