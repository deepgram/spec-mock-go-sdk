---
name: sdk-facade-conventions
description: Use when writing or reviewing Go code in pkg/. Covers package layout, deref helpers, pointer/value posture, naming idioms, type-switch patterns for streaming unions, and import aliasing. Route pipeline questions to sdk-codegen-flow, example/README rules to sdk-agentic-readiness, breaking changes to sdk-breaking-ceremony.
---

# Facade conventions

`pkg/` is the idiomatic Go facade over `api/`. Its job is to keep
customer-facing types and call surfaces stable across spec revisions.
These conventions are how it does that consistently.

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

Per `sdk-agentic-readiness`, every `.go` file in `pkg/` should be
single-concept. Don't co-locate REST and WebSocket types in one file.

### Legacy: `pkg/api/`

`pkg/api/{product}/.../interfaces/` is **legacy bootstrap** inherited
from `deepgram-go-sdk`. The split (REST under `pkg/client/`, WebSocket
under `pkg/api/`) was an accident of history; nothing about the spec
pipeline requires it. The migration plan: as each product moves through
`deepgram/spec` → `deepgram/spec-codegen-go` → here, its
`pkg/api/{product}/` subtree gets retired in favour of
`pkg/client/{product}/{ver}/{transport}/`.

When you add new code, use the target layout. When you edit existing
code at `pkg/api/listen/...` (the only product partly plumbed today),
match the surrounding style — the wholesale move is deferred to its
own refactor to avoid mid-flight risk.

## Import aliases

Use these consistently. Never dot-import. Never reference generated types
without an explicit alias:

```go
import (
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	ws        "github.com/deepgram/spec-mock-go-sdk/api/transport/websocket"
	sm        "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
)
```

The alias makes the conversion direction unambiguous in every reader's
head: `spectypes.X` is always wire-shape, never customer-shape.

## Pointer/value posture

| Layer | Field posture |
|---|---|
| `api/types/` | All optional fields are pointers (`*string`, `*int32`, `*float32`, `*bool`). Generated. Do not edit. |
| `pkg/` | Primitives are values. Nested structs stay pointers if they were pointer in the generated layer. Slices and maps stay as-is. |

Customers must never need to nil-check a primitive field. If a
customer-facing field semantically requires three-state (present /
absent / zero) — like `Speaker *int` to distinguish "no speaker
detected" from "speaker 0" — keep the pointer and document why.

## Deref helpers

These canonical helpers live in `convert.go` for each transport package.
Reuse them. Add new ones in the same style only when a new primitive
shape appears:

```go
func derefStr(p *string)  string  { if p == nil { return "" }; return *p }
func derefBool(p *bool)   bool    { if p == nil { return false }; return *p }
func derefInt32(p *int32) int     { if p == nil { return 0 }; return int(*p) }
func derefInt64(p *int64) int     { if p == nil { return 0 }; return int(*p) }
func derefF32(p *float32) float64 { if p == nil { return 0 }; return float64(*p) }
func int64ToPtrInt(p *int64) *int                         // when nil-meaning matters
func f32ToPtrF64(p *float32) *float64
```

If a generated field widens in a regen (`*float32 → *float64`), drop the
`f32ToPtrF64` call and use the pointer directly. Delete the helper if no
other call site uses it. The Go compiler will tell you.

## Converters

The facade has converters going **both directions**. Each direction has
its own naming pattern and its own rule for handling additive changes.

### Output direction: `convert{GeneratedTypeName}`

Every customer-facing response type has a converter from the generated
shape. Top-level nil check returns nil. Inside, never nil-check
primitives — use the deref helpers.

```go
func convertTranscribeOutput(in *spectypes.TranscribeOutput) *PreRecordedResponse {
	if in == nil {
		return nil
	}
	return &PreRecordedResponse{
		RequestID: derefStr(in.RequestId), // generated camelCase → idiomatic ID
		Duration:  derefF32(in.Duration),  // pointer + widen f32 → f64
		// ...
	}
}
```

One converter per type pair. Name `convert{GeneratedTypeName}`. Place
all converters in `convert.go`; helpers at the top, type converters
below in dependency order (leaf types first).

### Input direction: `optionsTo{GeneratedInputType}`

Some operations have an idiomatic options struct (e.g.
`PreRecordedTranscriptionOptions`) that maps onto a generated input
type (e.g. `spectypes.TranscribeInput`). Those converters go the
opposite direction — customer value-types → generated pointer-types —
and follow a different naming pattern:

```go
func optionsToTranscribeInput(o *interfaces.PreRecordedTranscriptionOptions) *spectypes.TranscribeInput {
	in := &spectypes.TranscribeInput{}
	if o == nil {
		return in
	}
	if o.Callback != "" {
		v := o.Callback
		in.Callback = &v
	}
	if o.Channels != 0 {
		v := int32(o.Channels)
		in.Channels = &v
	}
	if len(o.Keywords) > 0 {
		in.Keywords = o.Keywords
	}
	// ...
	return in
}
```

One converter per `(options struct → generated input type)` pair.
Name `optionsTo{GeneratedInputType}`. Each facade field gets a
non-zero check + pointer allocation; zero-valued facade fields leave
the generated field as `nil` so they don't surface in the wire query
string. This is the canonical input-side wiring shape, and the LLM
regen pipeline looks for exactly this pattern when extending the
facade for additive spec changes.

## Additive change rule

When a generated input or output type gains a field (`FIELD_ADDED`
classified as `absorbed-silently` by the regen pipeline), the facade
MUST be updated — silently, but definitely. "Silently" means no PR
comment, no `BREAKING_CHANGES.md` entry, no merge gate — the
customer-facing signature stays compatible. It does NOT mean "ignore
the change".

### What "extend the facade" means concretely

When `spectypes.{InputType}.{NewField}` appears:

1. **Ensure the facade options struct exposes a matching field.** Look
   at the options struct paired with this input type (e.g.
   `PreRecordedTranscriptionOptions` for `TranscribeInput`). If it
   already has the field with a sensible Go-idiomatic value type, leave
   it. If not, add it — same field name in Go-idiomatic case (e.g.
   `Numerals bool`, `Alternatives int`), positioned alphabetically.

2. **Extend the `optionsTo{InputType}` converter.** Add a non-zero check
   + pointer allocation block in alphabetical position. Match the
   existing style exactly (helpers / casts / value lifting). Never
   reorder existing blocks; insert in place.

3. **Update inline comments that list "dropped fields".** If the
   converter's docstring or inline comment enumerates fields the facade
   doesn't currently wire, REMOVE the now-wired field from that list.
   These lists must stay accurate or future runs will keep skipping the
   field.

### Dropped-list ↔ wiring consistency (mandatory check)

The dropped-fields list in `optionsTo{InputType}`'s docstring is a
contract with future runs. It enumerates facade-struct fields that
intentionally have NO wiring block in the converter (because the spec
doesn't model them yet). Two invariants MUST hold after every regen:

- **Every name in the dropped list has NO `if o.{Name}` block in the
  converter body.** Listing implies dropping.
- **Every facade-struct field NOT in the dropped list AND NOT a
  no-op-on-zero scalar (every field, in practice) MUST have an
  `if o.{Name}` block in the converter body.** Absence from the list
  implies wired.

Removing a name from the dropped list without adding the corresponding
wiring block is the most common regen bug. Never do that. If you find
yourself removing a name, the very next edit in the same response MUST
be the wiring block.

### Comprehensive scan rule (for multi-field syncs)

When the api/ diff contains MORE than one `FIELD_ADDED` on the same
input type — e.g., a big spec sync that adds ten or thirty fields at
once — walk every facade-options field (e.g. every field of
`PreRecordedTranscriptionOptions`) and ask:

1. Does this facade field have a matching spec field on the new
   `{InputType}`?
2. Does the converter already have an `if o.{Name}` block for it?

For every facade field where (1) is yes AND (2) is no, you owe a new
wiring block in this response. Do not stop after wiring a few obvious
ones; exhaustively cover the facade. A partial sync leaves silent
gaps where the customer can set a field that never reaches the wire.

When `spectypes.{OutputType}.{NewField}` appears:

1. **Ensure the customer response type exposes a matching field.** The
   value-type version with idiomatic case + JSON tag.

2. **Extend the `convert{OutputType}` converter.** Add the deref/copy
   line in alphabetical position alongside its peers.

3. **Emit an `Example_*` test if a new public symbol was added.** See
   `sdk-agentic-readiness`.

### Why this is non-breaking

The facade is designed so that additive changes never break customer
code:

- **Options structs are passed by pointer with named-field literals.**
  `&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3"}` keeps
  compiling when new fields appear. Existing callers don't need to
  set the new field.
- **Response structs use value-type primitives.** Zero is the absence
  signal; customers don't nil-check.
- **Method signatures take options-by-pointer, not options-by-value.**
  `FromURL(ctx, url, opts *Options)` keeps compiling when `Options`
  grows. New convenience constructors / helpers can be added without
  touching existing ones.
- **Converters use the deref helpers**, which already handle nil
  pointers. So a generated field that's optional in the wire shape
  surfaces in the facade as a zero-valued primitive — no nil checks
  needed in customer code.

When you add a new field to the facade in response to a generated
field, you are extending a stable, additive-friendly shape. You are
not changing a contract.

## Naming idioms

Go's golden rule for initialisms: caps stay caps. The generator emits
`RequestId`, `ApiVersion`, `DgRequestId`. The facade exposes
`RequestID`, `APIVersion`, `DgRequestID`. Apply this consistently.

| Generated | Idiomatic Go |
|---|---|
| `RequestId` | `RequestID` |
| `ApiVersion` | `APIVersion` |
| `DgRequestId` | `DgRequestID` |
| `Url` | `URL` |
| `Json` | `JSON` |
| `Http` | `HTTP` |
| `Sha256` | `SHA256` |

The customer-facing JSON tag still uses the wire name
(`json:"request_id"` etc.) so on-wire compatibility is preserved.

## Enum / sentinel conversion

Generated enums are `type Sentiment string`. Customer types use `*string`
where empty means absent:

```go
func sentimentToPtrStr(s spectypes.Sentiment) *string {
	if s == "" {
		return nil
	}
	v := string(s)
	return &v
}
```

When a new variant is added in the generated layer, the converter
usually needs no change — the underlying type is already `string`. Only
edit if there's a customer-facing switch on the value.

## Server-stream type switches

WebSocket converters use a type switch on the generated union:

```go
func ConvertServerMessage(m ws.ServerStream) any {
	switch v := m.(type) {
	case *ws.ServerStreamMemberResults:        return convertResults(&v.Value)
	case *ws.ServerStreamMemberMetadata:       return convertMetadata(&v.Value)
	case *ws.ServerStreamMemberSpeechStarted:  return convertSpeechStarted(&v.Value)
	// new variants: add a case + a customer-side event type
	default:                                   return nil
	}
}
```

When a new server message variant appears in `api/`, extend this switch
AND emit a new customer-side struct in the WebSocket facade's
`interfaces/types.go`. Wire it through the router and handler-call code
path the same way the existing events are wired. (Today that interfaces
package still lives at `pkg/client/listen/v1/websocket/interfaces/`; the
target home is `pkg/client/listen/v1/websocket/interfaces/` — see
"Legacy: pkg/api/" above.)

## When in doubt

- Match the surrounding style. The facade is small and consistent; new
  code should be indistinguishable from old.
- Prefer adding a helper to ad-hoc inline conversion.
- Comments explain WHY a conversion exists, not WHAT it does — e.g.
  `// f32 → f64 widens for printability` rather than `// convert float`.
- The agentic retrieval rule (one concept per file) takes precedence
  over locality. Split when in doubt.

## Related skills

- [`sdk-codegen-flow`](../sdk-codegen-flow/SKILL.md) — where this fits
  in the pipeline.
- [`sdk-agentic-readiness`](../sdk-agentic-readiness/SKILL.md) — every
  exported type/func above must ship with an Example_*.
- [`sdk-breaking-ceremony`](../sdk-breaking-ceremony/SKILL.md) — what
  to do if a spec change breaks customer signatures.
