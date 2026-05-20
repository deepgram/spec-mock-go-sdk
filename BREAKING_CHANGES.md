# Breaking Changes

This file is the cumulative record for tier 1 and tier 2 regen changes.
See `REGEN.md` section 15 for the full absorption rules used during codegen review.

Tier 0 changes are silently absorbed: new wire fields are added, wired,
and covered by tests without an entry here.

Tier 1 changes are absorbed with ceremony: removed wire fields stay public
but deprecated; converter wiring is removed and entries go here.

Tier 2 changes are unavoidable breaks: removed wire and public options fields
both disappear, and entries go here with review label context.

## No tier 1 or tier 2 changes have landed yet.
