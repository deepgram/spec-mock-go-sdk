# ADR-0002: Direct dependency choices + SBOM strategy

- Status: accepted
- Date: 2026-05-19
- Decider: DX team
- Supersedes: none

## Context

A customer-facing SDK should treat its dependency surface as carefully as its public-API surface. Every direct dependency is something a customer's build will pull in and something we're committing to keep current against CVE drift. The spec-driven prototype currently has 3 direct deps; this ADR records why each is there, what the risk is, and how the SBOM gets produced.

## Direct dependencies (as of 2026-05-19)

### `github.com/dvonthenen/websocket` v1.5.1-dyv.2

- **What:** WebSocket client/server library. Fork of `gorilla/websocket`.
- **Why we use it:** Carries the upstream `gorilla/websocket` API surface but with additional fixes / features the legacy `deepgram-go-sdk` relied on.
- **Risk level:** High. Single-maintainer fork of an upstream library that itself transitioned through community-maintenance status. Long-term support burden is on us if the maintainer goes inactive.
- **Mitigation plan:** Track upstream `gorilla/websocket` (or its current community fork) and evaluate migration each minor. Document any divergence we actually depend on so the migration path is well understood. Listed as a Track D1 (WS production-resilience) follow-up.

### `github.com/aws/aws-sdk-go-v2/service/sagemakerruntime` v1.39.8

- **What:** AWS SDK v2 client for SageMaker runtime, used by the SageMaker transport in `api/transport/sagemaker/`.
- **Why we use it:** SageMaker is one of the spec's `@supportsTransports` targets. AWS publishes this client themselves; rolling our own would be wrong.
- **Risk level:** Low. AWS-maintained, conventional minor cadence, no known supply-chain concerns.
- **Mitigation plan:** Track AWS-SDK-Go-v2 advisory feed via Dependabot once enabled.

### `github.com/aws/smithy-go` v1.25.1

- **What:** Runtime support library for Smithy-generated Go code (encoding helpers, middleware, document types).
- **Why we use it:** Required by the api/types/ output our spec-codegen-go emits. We don't control the version directly; it's pulled by the codegen output.
- **Risk level:** Low. AWS-maintained, in lockstep with smithy-go-codegen which we already depend on at build time.
- **Mitigation plan:** Pinned via go.mod. Bumped when spec-codegen-go bumps its smithy-go-codegen dep.

## Indirect dependencies

4 currently — all transitive from `aws-sdk-go-v2`. Listed in `go.sum`. Not separately audited because the AWS SDK owns their lifecycle.

## SBOM generation

### Format

CycloneDX 1.5 JSON. Industry standard, widely consumed by downstream supply-chain tooling.

### Tool

[`cyclonedx-gomod`](https://github.com/CycloneDX/cyclonedx-gomod) — generates SBOM from `go.mod` + `go.sum`. Maintained by the CycloneDX project, no controversial deps of its own.

### When

On every release. The release-please tag publish triggers `.github/workflows/sbom.yml`, which:

1. Checks out the tagged commit.
2. Installs `cyclonedx-gomod` from its release binary.
3. Generates `sbom.cdx.json` from the module graph.
4. Attaches the file to the GitHub release as a release artifact.

### Not in scope (yet)

- **Continuous SBOM publishing** (every commit) — overkill for an alpha SDK. Per-release coverage is sufficient signal.
- **Signing the SBOM** — useful once we have a customer audit story. Tracked as future work.
- **Vulnerability scanning** (Dependabot, govulncheck CI step) — separate concern, tracked separately.

## Adding a new direct dependency

Open a PR that:

1. Adds the require entry to `go.mod`.
2. Updates this ADR with a new section under "Direct dependencies" using the same shape (what / why / risk / mitigation).
3. Bumps the date on this ADR's header.

A new direct dependency without a matching ADR section MUST be rejected at review.

## Removing a direct dependency

Update this ADR. Move the section to a "## Removed dependencies" trailer with the removal date and reason. Don't delete the rationale — future regression debate references it.
