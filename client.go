package deepgram

// Client is the hand-written user-facing facade.
//
// It wraps generated types in the api package to provide an ergonomic surface
// that is stable across spec regeneration. The Smithy spec drives the api
// package contents; this struct and its methods are owned by humans and only
// change when we deliberately change the public Go SDK shape.
type Client struct {
	apiKey string
}

// NewClient returns a Client authenticated with the given Deepgram API key.
func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
