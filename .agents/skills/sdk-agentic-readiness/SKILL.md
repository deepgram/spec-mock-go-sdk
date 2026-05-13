---
name: sdk-agentic-readiness
description: Use when writing or reviewing any public API surface in this repo. The rules for scoring well on agentic retrieval tools (Context7, Fern agent-score, Mintlify, DeepScore) — runnable doc examples next to source, CI-validated examples, single-concept files, keyword-dense README opener, and a curated llms.txt. Sourced from dx-stack/docs/sdk-agentic-retrieval.md.
---

# Agentic readiness

This skill ports
[`deepgram/dx-stack/docs/sdk-agentic-retrieval.md`](https://github.com/deepgram/dx-stack/blob/main/docs/sdk-agentic-retrieval.md)
into this repo. That doc is the canonical strategy; this skill is the
Go-specific implementation guide for it.

## The pattern (from dx-stack)

> High-scoring libraries on agentic retrieval tools share three traits.
>
> **Examples live next to the code, not in a docs site.** Public APIs
> carry doc-comments with runnable example blocks. The retrieval engine
> extracts these as snippets — one method, one paragraph, one example.
>
> **Examples can't go stale.** They're either type-checked, executed as
> tests, or generated from real test fixtures. Drift is impossible by
> construction.
>
> **Files are single-concept and semantically named.** The retrieval
> engine's dedup step rewards non-overlapping content. Sprawling files
> get collapsed; tight files survive.

Treat these as hard requirements for `pkg/`. Treat them as nice-to-have
for `api/` (which is generated and not the surface customers should read
anyway).

## The five objectives

From the strategy doc, ported with the Go-specific implementation
clarified for each:

### 1. Every public API surface ships with a doc-comment containing at least one runnable example block

The Go-idiomatic format is **godoc + `Example_*` test functions**. The
two together satisfy the requirement:

```go
// PreRecordedResponse is the customer-facing response type returned by
// PreRecordedClient.FromURL, FromFile, and FromStream. Fields with
// pointer types are present only when the upstream API returned them;
// fields with value types are always populated (zero-valued if absent).
type PreRecordedResponse struct {
	// ...
}
```

```go
// pkg/client/listen/v1/rest/example_test.go
package restv1_test

import (
	"context"
	"fmt"
	"log"

	rest "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

// Example shows the minimal happy path: transcribe a URL and print the
// transcript of the first alternative.
func ExamplePreRecordedClient_FromURL() {
	client, err := rest.NewClient("DG_API_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.FromURL(context.Background(),
		"https://dpgr.am/spacewalk.wav", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
}
```

Why this works:

- `go doc` displays the example next to the type/function it documents.
- `pkg.go.dev` surfaces it on the type's page.
- `go test` runs every `Example_*` function. If the body doesn't
  compile, the test suite fails. If you add `// Output: ...` at the end
  of an Example, `go test` also asserts the printed output matches.
- The example lives at a known `file_test.go:line` next to the source it
  documents. Retrieval engines pick it up without a separate
  documentation site.

### 2. Examples are validated in CI

This is automatic in Go. `go test ./...` runs every `Example_*`. The
project's CI workflow (and `spec-idiomatic`'s build verifier) already
runs `go test ./...` on every regen, so:

- Examples that stop compiling break the build.
- Examples with `// Output: ...` assertions that fail break the build.
- A regen that changes a public signature must update the corresponding
  Example or it ships broken.

Add `// Output: ...` to examples where the output is deterministic.
Omit it where the output depends on a live API call or randomness;
`go test` will still compile-check those.

### 3. One concept per file

Split aggressively. Examples of correct grain:

- `pkg/client/listen/v1/rest/prerecorded.go` — only `FromURL` /
  `FromFile` / `FromStream` and their helpers.
- `pkg/client/listen/v1/rest/response.go` — only the
  `PreRecordedResponse` shape and its nested types.
- `pkg/client/listen/v1/rest/convert.go` — only converters and deref
  helpers.

Examples of wrong grain (don't do this):

- Mixing REST and WebSocket types in one file.
- Mixing client construction and response shapes in one file.
- A `helpers.go` or `utils.go` grab-bag.

`spec-idiomatic` refuses to write multi-concept files. If you're adding
something that doesn't fit the existing files, create a new one.

### 4. README opening paragraph is keyword-dense and lists runtime support

The first 500 characters set the library's representation in retrieval
indexes. Open with a one-sentence keyword-dense description and a
runtime-support sentence. Move badges below the fold.

Pattern (matching `dx-stack/docs/js-ecosystem.md:17-23`):

```
# spec-mock-go-sdk

spec-mock-go-sdk is the Go consumer of the Deepgram Smithy spec
pipeline, a mock SDK used to validate the agentic regen flow: it ships
generated wire types in `api/` and an idiomatic Go facade in `pkg/`
that absorbs spec wobbles so customer call signatures stay stable.
Runs on Go 1.22+. Not for customer use — see deepgram-go-sdk for the
official SDK.

[badges below the fold]

## ...
```

Keyword density target: every term a retrieval engine might match for
this SDK ("Deepgram", "Go", "Smithy", "facade", "agentic regen",
"wire types", "generated") shows up in the first paragraph.

### 5. Hand-curated llms.txt

`llms.txt` at repo root follows the [llmstxt.org](https://llmstxt.org/)
spec. It points retrieval engines at the canonical example files for
the SDK's major use cases.

```markdown
# spec-mock-go-sdk

> Go consumer of the Deepgram Smithy spec pipeline. Mock SDK for
> validating the agentic regen flow.

## Examples
- [Prerecorded transcription](pkg/client/listen/v1/rest/example_test.go): ExamplePreRecordedClient_FromURL
- [Live WebSocket](pkg/api/listen/v1/websocket/interfaces/example_test.go): ExampleConvertServerMessage
- ...
```

When you add a canonical example, add the line to `llms.txt`. When you
remove one, remove the line. `spec-idiomatic` maintains this on regen
but humans editing examples directly are responsible too.

## What this skill won't do

> Volume isn't the bottleneck. Concentration of correct, current,
> single-concept examples close to the source is.

So:

- Don't write Example_* you don't need.
- Don't add llms.txt entries for non-canonical examples.
- Don't pad the README opener to inflate keyword density — write the
  truth densely.

## How spec-idiomatic helps

When the regen runner edits or adds a public surface, it also:

- Generates or updates the corresponding `Example_*` function.
- Updates the README opening paragraph if the surface changed enough to
  affect the keyword density.
- Updates `llms.txt` if a canonical example was added, removed, or
  renamed.
- Refuses to write multi-concept files; splits into single-concept
  files if necessary.

Reviewer responsibility (per `sdk-pr-review`) is to verify these
artifacts moved together with the code.

## How to measure

Per dx-stack: behavioural benchmark on every release, comparable across
SDKs. `deepscore` is the planned tool. Week-over-week deltas tracked in
`dx-stack/docs/metrics.md`. Don't grade plumbing; grade outcomes.

## Related skills

- [`sdk-codegen-flow`](../sdk-codegen-flow/SKILL.md)
- [`sdk-facade-conventions`](../sdk-facade-conventions/SKILL.md)
- [`sdk-pr-review`](../sdk-pr-review/SKILL.md) — checks that examples,
  README, and llms.txt moved with the code.
