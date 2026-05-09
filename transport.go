package deepgram

import "context"

// Transport opens bidi-streaming sessions on behalf of generated clients.
//
// Built-in transports (transport/websocket, future transport/webrtc) live
// as subpackages of this module. Heavy transports that wrap a vendor SDK
// (SageMaker, future Vertex / Triton / Azure ML) live as separate Go
// modules the customer adds explicitly via go.mod. See ADR-0004.
//
// This minimal interface is a mock; the real one will land with the
// codegen-go runtime and will include frame send/receive primitives. For
// now Transport exists to demonstrate the pluggability surface.
type Transport interface {
	// Name returns the transport identifier matching the @supportsTransports
	// entry in the spec (for example "websocket" or "sagemaker").
	Name() string

	// Open establishes a streaming session against the given endpoint URL.
	// Real transports establish a connection; the mock returns nil.
	Open(ctx context.Context, url string) error
}
