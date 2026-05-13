# spec-mock-go-sdk

`spec-mock-go-sdk` is a Go consumer of the Deepgram Smithy spec pipeline. Its
job is to validate the regen flow end-to-end: every change to the spec lands
here as generated code in `api/`, and the idiomatic facade in `pkg/` absorbs
the wobble so customer-facing signatures stay stable.

This is **not** the official Deepgram Go SDK. It is the mock consumer the DX
team uses to develop and stress-test the agentic regen system (spec →
spec-codegen-go → spec-mock-go-sdk, with spec-idiomatic auto-regenerating the
facade on every spec change).

## Repo layout

| Path | Role | Edit policy |
|---|---|---|
| `api/` | Machine-generated wire types, transports, errors. Source of truth for the wire format. | **DO NOT EDIT.** Regenerated from spec every codegen run. Any edit will be wiped on the next regen. |
| `pkg/` | Idiomatic Go facade. Customer-facing types and call surfaces. Converts pointer-heavy generated types into value-heavy idiomatic types. | Edit freely. The facade exists to absorb api/ wobbles. |
| `examples/` | Runnable code samples per transport / use case. | Edit freely. Examples should compile and pass `go test ./examples/...`. |
| `tests/` | Unit and integration tests. | Edit freely. |
| `.agents/skills/` | Maintainer-facing skills (this directory). | Edit when conventions change. |
| `llms.txt` | Hand-curated index of canonical examples and docs for agentic retrieval tools. | Edit when public surface changes. |

## How a change lands

1. Someone merges a change in `deepgram/spec` (the Smithy IDL).
2. The fan-out workflow there opens a PR in this repo with `api/` regenerated.
3. `spec-idiomatic` fires on that PR, regenerates the `pkg/` facade, runs
   `go build ./...`, posts a work-done checklist comment, and pushes any
   edits back to the PR branch.
4. Human reviewer checks the diff per `sdk-pr-review`, acks any breaking
   ceremony per `sdk-breaking-ceremony`, merges.

## Skills

Skills live under `.agents/skills/<name>/SKILL.md`. Load whichever fits the
task you're about to do.

| Skill | Use when |
|---|---|
| [`sdk-codegen-flow`](.agents/skills/sdk-codegen-flow/SKILL.md) | You need to understand the api/ ↔ pkg/ split, the regen pipeline, or what `spec-idiomatic` does to a PR. |
| [`sdk-facade-conventions`](.agents/skills/sdk-facade-conventions/SKILL.md) | You're writing or reviewing Go code in `pkg/`. Deref helpers, pointer/value posture, naming idioms, type-switch patterns. |
| [`sdk-agentic-readiness`](.agents/skills/sdk-agentic-readiness/SKILL.md) | You're writing or reviewing any public API surface. The Example_*test, README-opener, single-concept-file, and llms.txt rules that keep this SDK scorable on retrieval benchmarks. |
| [`sdk-local-regen`](.agents/skills/sdk-local-regen/SKILL.md) | You want to run `spec-idiomatic` on your laptop against a synthetic api/ change before pushing anything to CI. |
| [`sdk-pr-review`](.agents/skills/sdk-pr-review/SKILL.md) | You're reviewing an auto-regen PR. Checklist for verifying the bot did the right thing. |
| [`sdk-breaking-ceremony`](.agents/skills/sdk-breaking-ceremony/SKILL.md) | The PR has a `regen/breaking-*` label or you're trying to decide what counts as breaking. The 3-tier model and reviewer playbook. |

## Cross-repo context

This repo is one of four:

- [`deepgram/spec`](https://github.com/deepgram/spec) — Smithy IDL source of truth.
- [`deepgram/spec-codegen-go`](https://github.com/deepgram/spec-codegen-go) — Java SmithyBuildPlugin that emits `api/`.
- [`deepgram/spec-mock-go-sdk`](https://github.com/deepgram/spec-mock-go-sdk) — this repo.
- [`deepgram/spec-idiomatic`](https://github.com/deepgram/spec-idiomatic) — agentic facade regenerator.

## Substrate philosophy

The **source code is the substrate**. `Example_*` functions, the README
opening paragraph, and `llms.txt` are **derived** from it. When the
substrate changes, the derived layer must be regenerated to match — that is
`spec-idiomatic`'s job, not yours, but you are responsible for verifying it
in review.

Prose that contradicts the substrate is wrong by definition. If you find a
mismatch, fix the prose, not the code.
