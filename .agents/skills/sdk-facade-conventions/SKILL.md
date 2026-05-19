---
name: sdk-facade-conventions
description: Use when editing or reviewing Go code in pkg/ in THIS repo. Repo-specific layout notes (where listen lives today, the pkg/api legacy split, which transports are wired) plus pointers to where the universal facade rules are documented. Route facade-author rules (3-tier model, converter naming, additive rule, wire-test contract) to spec-idiomatic's system.md and language/go.md.
---

# Facade conventions (repo-local)

The universal facade-author rules — 3-tier model, `convert<X>` /
`optionsTo<X>` patterns, additive rule, facade-mirror rule, wire-test
contract, BREAKING_CHANGES.md ownership — live in spec-idiomatic's
[`prompts/system.md`](https://github.com/deepgram/spec-idiomatic/blob/main/prompts/system.md)
and the Go syntax baselines (deref helpers, naming idioms, server-stream
type switches) live in
[`prompts/language/go.md`](https://github.com/deepgram/spec-idiomatic/blob/main/prompts/language/go.md).

This skill only documents what's specific to **this** repo.

## Package layout

Target layout — every transport for every product lives under
`pkg/client/{product}/{ver}/{transport}/`:

```
pkg/
└── client/{product}/{ver}/{transport}/   # REST, WebSocket, SageMaker, WebRTC
    ├── client.go                         # constructor + entry methods
    ├── {operation}.go                    # one file per RPC (prerecorded.go etc.)
    ├── response.go                       # customer-facing response value types
    ├── convert.go                        # generated → customer converters + helpers
    ├── types.go                          # customer-facing request types + options
    ├── constants.go                      # customer-visible enum values
    └── interfaces/                       # streaming-only sub-package
        ├── interfaces.go                 # handler interfaces customers implement
        ├── types.go                      # customer-facing event value types
        ├── convert.go                    # generated event → customer event converter
        └── constants.go                  # customer-visible message-type strings
```

Per [`sdk-agentic-readiness`](../sdk-agentic-readiness/SKILL.md), every
`.go` file in `pkg/` should be single-concept. Don't co-locate REST and
WebSocket types in one file.

## Legacy: `pkg/api/`

`pkg/api/{product}/.../interfaces/` is legacy bootstrap inherited from
`deepgram-go-sdk`. The split (REST under `pkg/client/`, WebSocket under
`pkg/api/`) was an accident of history; nothing about the spec pipeline
requires it. The migration plan: as each product moves through
`deepgram/spec` → `deepgram/spec-codegen-go` → here, its
`pkg/api/{product}/` subtree gets retired in favour of
`pkg/client/{product}/{ver}/{transport}/`.

When you add new code, use the target layout. When you edit existing
code at `pkg/api/listen/...`, match the surrounding style — the
wholesale move is deferred to its own refactor.

## What's wired through the new pipeline today

| Product | REST | WebSocket | Status |
|---|---|---|---|
| listen (transcribe) | ✅ | ⚠️ partial | Only product wired through `deepgram/spec` + `spec-codegen-go` + `spec-idiomatic` |
| speak | — | — | Legacy from `deepgram-go-sdk`, not yet plumbed |
| agent | — | — | Legacy |
| manage / projects | — | — | Legacy |

When a regen PR touches `api/**`, `spec-idiomatic` only modifies the
facade for products whose api/ shapes changed — which today means
`pkg/client/listen/v1/rest/` and `pkg/client/listen/v1/websocket/`.
Legacy `pkg/api/{other-product}/` subtrees stay frozen until they
migrate.

## Wire tests at `pkg/client/listen/v1/rest/wire_test.go`

This is the contract floor for the REST listen transport. One
`TestWires_<Field>` per facade-options field that flows through
`optionsToTranscribeInput`. Universal rules (when to add, when to
remove) are in spec-idiomatic's `system.md`; the file's own header
docstring spells out the per-field rule.

The build verifier runs `go test -count=1 ./pkg/...` on every regen
attempt. Tests that fail block the auto-commit, so silent gaps in
wiring get caught automatically.

## Related skills

- [`sdk-codegen-flow`](../sdk-codegen-flow/SKILL.md) — where this fits
  in the pipeline.
- [`sdk-agentic-readiness`](../sdk-agentic-readiness/SKILL.md) — every
  exported type/func above must ship with an `Example_*`.
- [`sdk-pr-review`](../sdk-pr-review/SKILL.md) — reviewer checklist for
  regen PRs.
