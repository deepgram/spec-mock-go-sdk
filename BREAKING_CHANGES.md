# Breaking changes

This file is authored by `spec-idiomatic`'s deterministic cleanup pass
on regen PRs. Entries describe customer-visible changes the facade
could not absorb silently — tier-2 unavoidable changes from
[`spec-idiomatic/prompts/system.md`](https://github.com/deepgram/spec-idiomatic/blob/main/prompts/system.md).

## Clean-slate facade reseed (2026-05-19)

- **Tier:** 2 (unavoidable-breaking)
- **Scope:** The entire `pkg/` tree was rewritten from scratch as a
  minimal facade aligned to current `api/`. The previous facade was
  inherited from `deepgram-go-sdk` and carried ~15k lines of legacy
  code unrelated to absorbing api/ wobbles. Post-reseed `pkg/` is
  ~1k lines.
- **Customer-visible impact:**
  - `pkg/client/listen/v1/rest/` is the only product wired today.
    Other products (`speak`, `agent`, `manage`, etc) are not present
    in `pkg/`. They will return per-product as each migrates through
    the spec pipeline.
  - Facade options structs are now co-located with their transport
    (e.g. `pkg/client/listen/v1/rest/options.go` exposes
    `PreRecordedTranscriptionOptions`). The previous shared
    `pkg/client/interfaces/v1/` package is gone; imports must be
    updated.
  - All `@internal`-tagged stem parameters that the spec excludes
    from the public surface are absent from the facade. Customer
    code that referenced `Alternatives`, `Channels`, `SampleRate`,
    or any other undocumented stem parameter on
    `PreRecordedTranscriptionOptions` will fail to compile.
- **Migration required:** Update imports to
  `github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest`
  and replace any reference to the legacy `pkg/api/` or
  `pkg/client/interfaces/` packages.
