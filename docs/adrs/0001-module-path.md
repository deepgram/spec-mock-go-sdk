# ADR-0001: Module path for the spec-driven Go SDK

- Status: accepted
- Date: 2026-05-19
- Decider: DX team
- Supersedes: none

## Context

The spec-driven SDK prototype lives in `github.com/deepgram/spec-mock-go-sdk` while it is being shaped. There is an existing official customer-facing Go SDK at `github.com/deepgram/deepgram-go-sdk` that is published, documented, and depended on. Once the spec-driven prototype is ready to ship, we need a module path strategy for the production SDK.

Two reasonable shapes:

1. **Replace in place.** Rename / repoint `deepgram/deepgram-go-sdk` to point at this work. Existing import paths keep working; existing customer code breaks at the major version. One go-getable module for customers to reach for.

2. **Side-by-side.** Keep `deepgram/deepgram-go-sdk` for the legacy SDK, publish the spec-driven SDK at a new path (e.g. `github.com/deepgram/sdk-go`). Two coexisting SDKs; customers opt in.

## Decision

Side-by-side. The spec-driven SDK ships at a NEW module path. The legacy SDK remains on its current import path indefinitely (or until separately deprecated and archived).

The new path is to be determined at first-release time, but the working assumption is `github.com/deepgram/sdk-go` (mirrors `deepgram/sdk-python`, `deepgram/sdk-js`, etc., once those ship with the same shape).

## Rationale

- **Customer pain is bounded.** Existing customers on `deepgram-go-sdk` are not forced to migrate. They migrate when they choose to. There is no "one day your build breaks because the SDK changed under you" event.
- **Alpha-stage discipline matches the path choice.** Per the prototype's working notes, "all interfaces can change without backwards compat" — that's exactly the wrong contract to ship on the import path customers already depend on. Side-by-side gives us the freedom to iterate without surprising anyone.
- **The migration story is opt-in.** A "should I move?" doc + a per-product surface comparison table lets customers make a deliberate decision rather than discovering it via a CI failure.
- **Future per-language SDKs converge on a clean naming pattern.** `sdk-go` / `sdk-python` / `sdk-js` is easier to communicate than "the new deepgram-go-sdk, which is different from the old deepgram-go-sdk".

## What we accept

- **Two SDKs coexist.** Documentation, examples, and support touchpoints have to acknowledge both. This is non-trivial overhead.
- **Some customer fragmentation.** Old code stays on old SDK longer. We must commit to legacy security/critical-bug coverage for a defined window.
- **Naming churn before first release is OK.** The path can still change before `v0.1.0` ships. It cannot change after.

## What we reject

- **Replace in place** is rejected because it inflicts a forced migration on every existing customer at major-version bump time, and the prototype's "interfaces can change without backwards compat" rhythm is incompatible with that customer-protection contract.
- **Hostile sunset of the legacy SDK** is rejected. There is no business reason to archive `deepgram-go-sdk` before the spec-driven SDK has earned customer trust.

## Open at decision time

- **Exact final import path.** Working assumption `sdk-go`. Finalised before first release.
- **Legacy SDK lifecycle.** When does `deepgram-go-sdk` go into security-only mode? When does it archive? Tracked separately.
- **Codemod.** A `gofmt`-style migration tool would lower the cost of side-by-side. Listed under D2 in the work plan but not blocking this decision.

## Consequences for other tracks

- **B2 (release-please).** Wires against `github.com/deepgram/spec-mock-go-sdk` for now; the rename to the final path is a future ADR that ships alongside the rename PR.
- **D2 (migration guide).** Lives in the new SDK's `docs/migration-from-deepgram-go-sdk.md`. The legacy SDK gets a one-line pointer.
- **D3 (internal pilot).** Internal team adopts the new module path directly, not via legacy aliasing.
