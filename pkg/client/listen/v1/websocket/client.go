// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import (
	"context"
	"errors"
	nethttp "net/http"
	"os"

	wstransport "github.com/deepgram/spec-mock-go-sdk/api/transport/websocket"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

const DefaultBaseURL = "wss://api.deepgram.com"

const streamPath = "/v1/listen"

// Client is the Listen WebSocket client. It holds credentials and
// transport configuration shared across all Connect calls.
//
// The zero value of Client is not usable. Construct via New or
// NewWithDefaults.
type Client struct {
	apiKey      string
	accessToken string
	baseURL     string
}

// New constructs a Client with explicit credentials. Either apiKey
// or accessToken must be non-empty for the dial to authenticate;
// accessToken takes precedence when both are set.
func New(apiKey, accessToken string) *Client {
	return &Client{
		apiKey:      apiKey,
		accessToken: accessToken,
		baseURL:     DefaultBaseURL,
	}
}

// NewWithDefaults reads DEEPGRAM_ACCESS_TOKEN and DEEPGRAM_API_KEY
// from the environment.
func NewWithDefaults() *Client {
	return New(
		os.Getenv("DEEPGRAM_API_KEY"),
		os.Getenv("DEEPGRAM_ACCESS_TOKEN"),
	)
}

// WithBaseURL returns a copy of the Client pointed at the given
// WebSocket base URL (e.g. "wss://api.deepgram.com").
func (c *Client) WithBaseURL(url string) *Client {
	out := *c
	out.baseURL = url
	return &out
}

// Connect opens a WebSocket stream to /v1/listen with the supplied
// options encoded as URL query parameters. The returned Stream is a
// bidirectional handle: callers send audio and control messages via
// Stream.SendAudio / Stream.Finalize / Stream.KeepAlive /
// Stream.CloseStream / Stream.Sync, and receive events via
// Stream.Recv. Close the Stream when done.
//
// The Configure variant of ClientStream (mid-session feature
// reconfiguration) is not exposed on Stream yet because its wire
// shape uses document.Interface; callers needing it can reach the
// underlying transport via api/transport/websocket directly.
func (c *Client) Connect(ctx context.Context, opts *LiveTranscriptionOptions) (*Stream, error) {
	headers, err := c.authHeaders()
	if err != nil {
		return nil, err
	}

	wireOpts := optionsToStreamInput(opts)
	query := streamInputQueryString(wireOpts)
	dialURL := c.baseURL + streamPath
	if query != "" {
		dialURL += "?" + query
	}

	transport, err := wstransport.OpenStream[spectypes.ClientStream, spectypes.ServerStream](
		ctx,
		dialURL,
		headers,
		func(msg spectypes.ClientStream) ([]byte, bool, error) {
			return spectypes.MarshalClientStream(msg)
		},
		spectypes.UnmarshalServerStream,
	)
	if err != nil {
		return nil, err
	}
	return &Stream{transport: transport}, nil
}

func (c *Client) authHeaders() (nethttp.Header, error) {
	h := nethttp.Header{}
	switch {
	case c.accessToken != "":
		h.Set("Authorization", "Bearer "+c.accessToken)
	case c.apiKey != "":
		h.Set("Authorization", "Token "+c.apiKey)
	default:
		return nil, errors.New("listen websocket: no credentials; set DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN, or pass them to New(...)")
	}
	return h, nil
}

// Stream is a single live transcription session. Methods are
// goroutine-safe for sends (the underlying transport serializes
// writes); Recv is expected to be driven from a single goroutine.
type Stream struct {
	transport wstransport.Stream[spectypes.ClientStream, spectypes.ServerStream]
}

// SendAudio transmits a raw audio chunk as a binary WebSocket frame.
// The chunk's encoding must match the Encoding/SampleRate/Channels
// options used at Connect time. An empty chunk is treated by the
// server as equivalent to CloseStream.
func (s *Stream) SendAudio(data []byte) error {
	return s.transport.Send(&spectypes.ClientStreamMemberAudio{
		Value: spectypes.AudioFrame{Data: data},
	})
}

// CloseStream sends a graceful end-of-audio marker. The server
// flushes pending results and emits a final MetadataEvent before
// closing the WebSocket.
func (s *Stream) CloseStream() error {
	return s.transport.Send(&spectypes.ClientStreamMemberCloseStream{
		Value: spectypes.CloseStream{},
	})
}

// Finalize forces the server to emit a final transcript for the
// open utterance. Pass channel = -1 to finalize all channels.
func (s *Stream) Finalize(channel int) error {
	v := spectypes.Finalize{}
	if channel >= 0 {
		c := int32(channel)
		v.Channel = &c
	}
	return s.transport.Send(&spectypes.ClientStreamMemberFinalize{Value: v})
}

// KeepAlive resets client and driver inactivity timers. Send
// periodically when there's no audio to transmit but the session
// should stay open.
func (s *Stream) KeepAlive() error {
	return s.transport.Send(&spectypes.ClientStreamMemberKeepAlive{
		Value: spectypes.KeepAlive{},
	})
}

// Sync sends a client-side sync marker with the supplied id. The
// server echoes it back as a SyncEvent. Used to align test
// pipelines.
func (s *Stream) Sync(id int64) error {
	return s.transport.Send(&spectypes.ClientStreamMemberSync{
		Value: spectypes.ClientSync{Id: &id},
	})
}

// Recv blocks until the next event arrives or the connection
// closes. Returns nil for messages the SDK does not recognize so
// callers can continue the loop on forward-compat variants. Returns
// io.EOF (or other transport error) when the connection ends.
func (s *Stream) Recv() (Event, error) {
	msg, err := s.transport.Recv()
	if err != nil {
		return nil, err
	}
	return fromServerStream(msg), nil
}

// Close terminates the WebSocket connection. Idempotent; safe to
// call after CloseStream.
func (s *Stream) Close() error {
	return s.transport.Close()
}
