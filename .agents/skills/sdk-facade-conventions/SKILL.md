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

Every customer-facing type has a converter. Top-level nil check returns
nil. Inside, never nil-check primitives — use the deref helpers.

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
package still lives at `pkg/api/listen/v1/websocket/interfaces/`; the
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
