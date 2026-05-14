<!-- spec-idiomatic:work-done -->
## spec-idiomatic — work done

**Last run:** ${LAST_RUN}

### Steps

- [x] Materialise before/after api/ snapshots from PR base SHA
- [x] Diff api/ and classify changes (silently-absorbed / absorbed-with-ceremony / unavoidable)
- [x] Parse pkg/ Go facade AST via the shelled-out Go AST helper
- [x] Load consumer's `AGENTS.md` + every `.agents/skills/*/SKILL.md`
- [x] Compose orchestrator prompt + invoke `claude-opus-4-7`
- [x] Apply returned edits under pkg/ (refused any path containing `..`)
- [x] Write `BREAKING_CHANGES.md` for tier-1 / tier-2 changes the LLM flagged
- [x] Run `go test -count=1 ./pkg/...` (retried up to 3 times with build error appended)
- [x] Write structured JSON report (`spec-idiomatic-report.json`)
- [x] Apply `regen/breaking-absorbed` / `regen/breaking-unavoidable` labels per report
- [x] Edit tier-1 / tier-2 comments in place via HTML-sentinel lookup
- [x] Set `spec-idiomatic/breaking-acked` status check
- [x] Commit + push facade updates to PR branch (no-op when no edits)
- [x] Tick the `workflow-checkbox:idiomatic` checkbox in PR description

### Results

| Metric | Count |
|---|---|
| absorbed-with-ceremony | ${ABSORBED_COUNT} |
| unavoidable | ${UNAVOIDABLE_COUNT} |
| edits applied to pkg/ | ${EDITS_APPLIED} |
| build retries | ${ATTEMPTS_MINUS_ONE} |

### Per-change classification

${CHANGES_LIST}

This comment is edited in place on every workflow run — see the `Last run` timestamp above for the latest pass.
See [deepgram/spec-idiomatic](https://github.com/deepgram/spec-idiomatic) for runner internals.
