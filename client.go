package deepgram

// Client is the hand-written user-facing facade.
//
// It wraps generated types in the api package to provide an ergonomic surface
// that is stable across spec regeneration. The Smithy spec drives the api
// package contents; this struct and its methods are owned by humans and only
// change when we deliberately change the public Go SDK shape.
type Client struct {
	apiKey    string
	transport Transport
}

// Option configures a Client at construction time.
type Option func(*Client)

// WithTransport selects the bidi-streaming transport.
//
// Built-in transports (transport/websocket, future transport/webrtc) ship
// with this module. Heavy transports that wrap a vendor SDK (SageMaker,
// future Vertex / Triton / Azure ML) ship as separate Go modules the
// customer adds to go.mod explicitly.
func WithTransport(t Transport) Option {
	return func(c *Client) { c.transport = t }
}

// NewClient returns a Client authenticated with the given Deepgram API key.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{apiKey: apiKey}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
