# smithy-mock-go-sdk

Mock fan-out target for [`deepgram/spec`](https://github.com/deepgram/spec).

This repo simulates the role of `deepgram/deepgram-go-sdk` in the spec
release workflow: when `deepgram/spec` cuts a release, its fan-out workflow
runs smithy-go codegen and opens a PR here against a fixed branch
(`spec-sync`).

## Layout

```
api/          regen target — smithy-go output (NEVER edit by hand)
client.go     hand-written facade — stable across spec regen
doc.go        package documentation
go.mod        Go module definition
.github/      CI (go vet, go test, go build)
```

## Two-layer SDK pattern

This mock demonstrates the layered approach we want for real Deepgram SDKs:

1. **`api/`** — generated, deterministic, regenerates on every spec release.
   Treat as opaque except where the facade exposes it.
2. **`client.go`** (and any future hand-written files at the package root) —
   stable user-facing API. Wraps `api/` types in an ergonomic shell that
   survives regeneration.

When the spec adds a field, regen updates `api/`. The facade only needs
updating if that field becomes user-visible. When the spec adds a new
operation, regen updates `api/` AND the facade gets a new method.

## Status

Mock. Not consumed by anyone. Safe to break, safe to delete.
