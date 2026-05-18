---
name: sdk-breaking-ceremony
description: Use when a regen PR carries a `regen/breaking-*` label, when you're trying to decide whether a spec change should be silently absorbed or surfaced as breaking, or when you need to ack a tier-2 break to unblock merge. Covers the 3-tier model, the labels and status checks each tier triggers, and the reviewer playbook.
---

# Breaking-change ceremony

Most spec changes are absorbed silently by the facade. Some can't be
absorbed without losing semantic information. A few can't be absorbed
at all. Each class gets a different ceremony.

## The 3 tiers

| Tier | Label | Merge block | What it means |
|---|---|---|---|
| **0** | (none) | no | Silently absorbed. New optional field, new union variant added, idiomatic rename (`Id`→`ID`), type widening that preserves precision. Customer signatures stable, no human attention needed. |
| **1** | `regen/breaking-absorbed` | no | Absorbed with ceremony. Facade could absorb the change but had to drop semantic info. Field removed → customer field kept but always nil. Required → optional flip. Union variant removed → customer event handler unreachable. Reviewer reads the absorbed-breaking comment and decides whether the absorption is acceptable. |
| **2** | `regen/breaking-unavoidable` | yes | Can't absorb. Operation removed, type kind flipped, fundamental contract change. Customer code WILL break at compile or runtime. Status check `spec-idiomatic/breaking-acked` fails until a CODEOWNER applies the `regen/breaking-acked` label. |

## What triggers what

`spec-idiomatic`'s `DiffAnalyzer` classifies each change. The mapping
from change kind to tier:

| Change | Tier |
|---|---|
| `FIELD_ADDED` (new optional field) | 0 |
| `UNION_VARIANT_ADDED` (new server message variant) | 0 |
| `TYPE_ADDED` (new shape entirely) | 0 |
| `FIELD_RENAMED` (idiomatic-only, e.g. `Id`→`ID`) | 0 |
| `FIELD_TYPE_WIDENED` (e.g. `*float32`→`*float64`) | 0 |
| `FIELD_REMOVED` where the facade can default to zero/nil | 1 |
| Required → optional flip on an existing field | 1 |
| Optional → required flip on an existing field | 1 |
| `UNION_VARIANT_REMOVED` where customer event handler still compiles | 1 |
| `FIELD_TYPE_NARROWED` where a lossy convert exists | 1 |
| `OP_REMOVED` (operation gone from the spec) | 2 |
| `TYPE_KIND_FLIPPED` (struct ↔ enum, etc.) | 2 |
| `FIELD_REMOVED` where the customer surface depends on it irrecoverably | 2 |
| `FIELD_REMOVED` where the removal is intentional spec hygiene | 2 |
| `FIELD_TYPE_NARROWED` with no safe conversion | 2 |
| Any change the LLM explicitly declares unabsorbable | 2 |

The LLM may upgrade a change from tier 0 or 1 to tier 2 if it cannot
produce facade edits that compile cleanly. It cannot downgrade.

### Spec-hygiene removals are tier 2, not tier 1

The default for `FIELD_REMOVED` is tier 1: keep the customer field on
the facade struct, drop the wiring in `optionsTo{InputType}`, write a
`TestDropped_<Field>` to lock in the absorption. That's the right move
for accidental or hot-fix removals where customer source-compat is
worth more than wire fidelity.

It is **not** the right move when the removal is intentional spec
hygiene — a field being pulled out of the public API because it should
never have been exposed in the first place. Indicators that a removal
is hygiene rather than accident:

- The triggering spec PR body mentions `@internal` tagging, deprecation
  cleanup, or "fields that shouldn't be public" / "audit" / "stem-only".
- The removed field corresponds to a stem parameter explicitly marked
  "not publicly documented" or "for internal use only" in
  `stem/src/handlers/queries.rs`.
- The removed field is absent from `developers.deepgram.com`.

For hygiene removals, upgrade the classification to tier 2:

- Remove the member from the facade `*Options` struct entirely.
- Remove the corresponding `TestWires_<Field>` / `TestDropped_<Field>`
  test (the field no longer exists on the facade, so neither test
  applies).
- Emit a tier-2 `breaking_changes[]` entry with `reason` explaining
  the hygiene justification (e.g. "internal-only stem parameter,
  removed from public spec; facade field dropped to surface the break
  to customers rather than silently no-op the setter").
- The runner labels the PR `regen/breaking-unavoidable` and the
  `spec-idiomatic/breaking-acked` status check fails until a CODEOWNER
  acks. That gate is the correct ceremony: customer code that referenced
  the field will fail to compile on the next SDK upgrade, exactly as
  the audit intends.

Customers using docs already had no reason to set these fields, so the
compile break only catches customers reaching into undocumented stem
parameters — exactly the population that should be flagged.

## Ceremony per tier

### Tier 0

Bot does its thing, commits the facade, posts the work-done checklist
comment. No labels. No special review attention beyond the standard
`sdk-pr-review` checklist.

### Tier 1 (absorbed-breaking)

Bot does:

1. Applies label `regen/breaking-absorbed` to the PR.
2. Adds an entry to `BREAKING_CHANGES.md` at repo root describing
   the shape, member, what got nil/zero/dropped, and what customer
   code that read the lost data will now see.
3. Posts a dedicated PR comment titled **"⚠️ Absorbed breaking
   changes"** with the same content in human-readable form.
4. Commits and pushes facade edits as normal.

Reviewer must:

1. Read the bot's comment and the `BREAKING_CHANGES.md` entry.
2. Decide whether silent absorption is acceptable for this surface.
3. If acceptable: merge as normal. The label travels with the PR
   into git history as the audit trail.
4. If not acceptable: apply the `regen/override-to-unavoidable`
   label, comment with reasoning, and the PR is treated as tier 2
   from that point forward. (The next regen will produce different
   absorption attempts based on the override signal in the spec
   PR.)

### Tier 2 (unavoidable-breaking)

Bot does:

1. Applies label `regen/breaking-unavoidable` to the PR.
2. Adds an entry to `BREAKING_CHANGES.md` at repo root with a
   prominent **Migration required** section.
3. Posts a dedicated PR comment titled **"🛑 Unavoidable breaking
   changes"** with the migration guidance.
4. Sets the `spec-idiomatic/breaking-acked` status check to FAILURE.
5. Commits and pushes whatever partial facade edits it could
   produce (may not compile cleanly; CI will reflect this).

Reviewer (must be CODEOWNER):

1. Read the bot's comment and the `BREAKING_CHANGES.md` entry.
2. Decide if the break is acceptable. Common reasons it is:
   - The removed operation is deprecated upstream and customers
     were already warned.
   - The narrowed type matches a security or correctness fix.
   - The kind flip is part of a major version bump.
3. If acceptable: apply the `regen/breaking-acked` label. The
   status check re-evaluates and turns green. Merge.
4. If not acceptable: do not merge. The PR sits open until the
   upstream spec change is reverted, redesigned, or scheduled
   for a major version bump in `deepgram/spec`.

## BREAKING_CHANGES.md format

The bot writes this file. Don't edit it by hand on a regen PR —
your edits will be overwritten on the next bot run. If you need to
add context, comment on the PR; the human comment is the audit
trail alongside the file.

Format (per entry):

```markdown
## TranscribeOutput.Metadata removed

- **Tier:** 1 (absorbed-breaking)
- **Shape:** `spectypes.TranscribeOutput` (generated, in `api/types`)
- **Member:** `Metadata *ResponseMetadata`
- **Customer-visible impact:** Customer-facing
  `PreRecordedResponse.Metadata *Metadata` field is preserved in the
  facade to keep the Go API stable, but will now always be nil. Any
  customer code that read `resp.Metadata.RequestID`, model info,
  warnings, etc., will nil-dereference unless it nil-checks first.
- **Facade behaviour:** `convertTranscribeOutput` no longer references
  `in.Metadata`; hardcodes `Metadata: nil`. The
  `convertResponseMetadata` helper has been removed.
```

For tier 2 entries, add a **Migration required** section with the
specific customer code changes needed.

## When to upgrade or downgrade by hand

You can override the bot's tier classification in two directions:

- **Upgrade tier 1 → tier 2** (treat an absorbed break as
  unavoidable): apply the `regen/override-to-unavoidable` label.
  Useful when the silent absorption is technically possible but
  semantically misleading enough that you want customers warned via
  a major version bump.
- **Downgrade is not allowed.** If the bot classified something as
  tier 2, it could not produce clean facade edits. Forcing it
  through as tier 1 would ship a broken SDK. The right move is to
  fix the underlying spec issue and let the regen run again.

## Status check details

The `spec-idiomatic/breaking-acked` status check is set by the bot
itself based on the PR's labels:

- No `regen/breaking-*` labels → check passes (skipped).
- `regen/breaking-absorbed` only → check passes (acked by virtue of
  being absorbable).
- `regen/breaking-unavoidable` and no `regen/breaking-acked` → check
  fails.
- `regen/breaking-unavoidable` and `regen/breaking-acked` → check
  passes.

The check re-evaluates on label changes. Applying `regen/breaking-acked`
turns the check green within seconds.

## Why the ceremony exists

A regen that silently breaks customer code is worse than a regen that
loudly breaks customer code. The labels and the status check are the
mechanism that ensures every break is a deliberate decision someone
ack'd, not an accident no one noticed until customer support tickets
arrived.

## Related skills

- [`sdk-codegen-flow`](../sdk-codegen-flow/SKILL.md) — where the bot
  does the classifying.
- [`sdk-pr-review`](../sdk-pr-review/SKILL.md) — the broader review
  checklist this ceremony fits into.
- [`sdk-facade-conventions`](../sdk-facade-conventions/SKILL.md) —
  what absorption actually looks like in code.
