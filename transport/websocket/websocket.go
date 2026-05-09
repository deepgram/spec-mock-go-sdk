// Package websocket is the built-in WebSocket transport, shipping as a
// subpackage of the main SDK. Customers do not install anything extra to
// use it. The real implementation will use gorilla/websocket or
// nhooyr.io/websocket; this mock leaves that choice to the codegen-go
// runtime.
package websocket

import (
	"context"

	deepgram "github.com/deepgram/smithy-mock-go-sdk"
)

// Transport is the WebSocket implementation of deepgram.Transport.
type Transport struct{}

// New returns a WebSocket transport with default configuration.
func New() deepgram.Transport {
	return &Transport{}
}

// Name returns "websocket".
func (t *Transport) Name() string { return "websocket" }

// Open establishes a WebSocket connection. Mock — no-op.
func (t *Transport) Open(ctx context.Context, url string) error { return nil }
