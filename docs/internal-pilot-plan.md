# Internal pilot plan

The single highest-signal validation for the spec-driven SDK is one internal team running it for a real workload. This document captures what that engagement needs and what success looks like, so it can be picked up by a human collaborator when scheduling allows.

## Why a pilot

Tests passing + lsp clean + build green ≠ done. `httptest` synthetics catch wire-mapping regressions, not the things that go wrong against a real running service: latency, retries, idle-connection handling, edge-case error responses, response payload shapes the codegen got slightly wrong.

The pilot is how those find us before they find a customer.

## Pilot candidates

Looking for one of the following profiles:

- **Internal team using `deepgram-go-sdk` today.** Lowest activation cost — they already have a Go codebase calling the legacy SDK. Switch the import path on a feature branch, see what breaks.
- **New internal Go integration.** Lower-effort migration, but they have no baseline to compare against. Still valuable.
- **An external example that DX maintains.** `deepgram/examples`, `deepgram/recipes`, or one of the per-product examples directories. Less interactive than a real workload but signal-rich.

## What the pilot needs from us

1. **Pinned alpha tag.** Once `v0.1.0-alpha.1` cuts via release-please, the pilot tracks a specific version. No `@latest` against an alpha.
2. **Support channel.** Slack DM or a dedicated `#dx-pilot-sdk-go` thread. Sub-day response time during the pilot window.
3. **Migration guide as the starting point.** [`docs/migration-from-deepgram-go-sdk.md`](./migration-from-deepgram-go-sdk.md) walks them through the rewrites.
4. **Bug ownership.** Anything they hit is OUR bug, not theirs. We fix it on the spec-driven side; they don't work around it.
5. **A short retro.** ~30 minutes at the end of the pilot window. Captures what worked, what didn't, what surprised them.

## What we need from the pilot

1. **One real workload.** Not a demo; a thing customers (or internal systems) actually depend on. Production traffic or staging traffic of equivalent shape.
2. **Honest reporting.** "This is weird" / "I don't get why it's named X" / "I had to write a workaround" — all valuable. Soft signals more than hard bugs.
3. **A migration log.** A diff of what changed in their integration, before/after. Surfaces the bulk of the rewrite cost.
4. **Bug filing in `deepgram/spec-mock-go-sdk` issues**, not Slack DMs. Slack for triage / quick questions; issues for tracked work.

## Success criteria

- Pilot completes the migration without abandoning it midway because something is unfixable.
- All bugs surfaced are addressed in the SDK (not papered over in the pilot's code).
- Retro identifies at least one architectural surprise — something we'd otherwise have only learned from a paying customer.
- The pilot's code stays on the spec-driven SDK after the pilot window ends.

## What this is NOT

- A beta program. We're not signing up multiple teams.
- A way to ship faster. We're not setting an aggressive timeline that compromises code quality. If the pilot stalls, that's a finding, not a failure of the pilot.
- A substitute for proper CI / wire-test coverage. The pilot is the last validation gate, not the only one.

## Status

**Not yet scheduled.** Blocking on:

- A pilot team identified. Track A and B work has to be visibly green before we ask anyone to commit time.
- `v0.1.0-alpha.1` tag cut. Cannot pilot against `main`.

Action: when the SDK is ready to pitch a pilot, post in `#team-devrel-internal` with this doc as the briefing and a sign-up window.
