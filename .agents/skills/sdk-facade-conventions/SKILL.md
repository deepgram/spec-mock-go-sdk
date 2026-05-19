---
name: sdk-facade-conventions
description: Use when editing or reviewing Go code in pkg/ in THIS repo. Repo-local layout notes (which products are wired, where wire tests live) plus pointers to where the universal facade rules are documented. Route facade-author rules (3-tier model, converter naming, additive rule, wire-test contract) to spec-idiomatic's system.md and language/go.md.
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

Every transport for every product lives under
`pkg/client/{product}/{ver}/{transport}/`, with options + response +
converters + client + wire tests co-located in a single package:

```
pkg/client/{product}/{ver}/{transport}/
├── client.go         # Client + Connect / FromURL / FromFile / FromStream
├── options.go        # <Operation>Options facade struct
├── response.go       # customer-facing response value types  (REST)
├── events.go         # customer-facing event types           (streaming)
├── convert.go        # deref helpers + optionsTo<X> + convert<X> / fromX
├── example_test.go   # Example_* smoke tests
├── wire_helpers.go   # requireWired / requireDropped helpers
└── wire_test.go      # one TestWires_<Field> per facade-options field
```

Per [`sdk-agentic-readiness`](../sdk-agentic-readiness/SKILL.md), every
`.go` file in `pkg/` should be single-concept.

## Products wired today

| Product | REST | WebSocket |
|---|---|---|
| listen (transcribe / live) | ✅ | ✅ |

Other products return per-product as each migrates into the spec
pipeline. When a regen PR touches `api/**`, `spec-idiomatic` only
modifies facades for products whose api/ shapes changed.

## Wire tests at `pkg/client/listen/v1/{rest,websocket}/wire_test.go`

These are the contract floor for each transport. One
`TestWires_<Field>` per facade-options field that flows through the
matching `optionsTo<X>` converter. Universal rules (when to add, when
to remove) are in spec-idiomatic's `system.md`; the file's own header
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
