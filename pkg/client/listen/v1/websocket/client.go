// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"net/url"
	"os"
	"time"

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
	config      *Config
}

// New constructs a Client with explicit credentials. Either apiKey
// or accessToken must be non-empty for the dial to authenticate;
// accessToken takes precedence when both are set.
//
// Example:
//
//	client := ws.New("your-api-key", "")
//	stream, err := client.Connect(ctx, opts)
//
// See [ExampleNew] for the full runnable form.
func New(apiKey, accessToken string) *Client {
	return &Client{
		apiKey:      apiKey,
		accessToken: accessToken,
		baseURL:     DefaultBaseURL,
	}
}

// NewWithDefaults reads DEEPGRAM_ACCESS_TOKEN and DEEPGRAM_API_KEY
// from the environment.
//
// Example:
//
//	client := ws.NewWithDefaults()
//	stream, err := client.Connect(ctx, &ws.LiveTranscriptionOptions{
//	    Model:      "nova-3",
//	    Encoding:   "linear16",
//	    SampleRate: 16000,
//	})
//
// See [ExampleNewWithDefaults] for the full runnable form.
func NewWithDefaults() *Client {
	return New(
		os.Getenv("DEEPGRAM_API_KEY"),
		os.Getenv("DEEPGRAM_ACCESS_TOKEN"),
	)
}

// WithBaseURL returns a copy of the Client pointed at the given
// WebSocket base URL (e.g. "wss://api.deepgram.com").
//
// Example:
//
//	client := ws.NewWithDefaults().WithBaseURL("wss://staging.api.deepgram.com")
//
// See [ExampleClient_WithBaseURL] for the full runnable form.
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
//
// Example:
//
//	stream, err := client.Connect(ctx, &ws.LiveTranscriptionOptions{
//	    Model:          "nova-3",
//	    Encoding:       "linear16",
//	    SampleRate:     16000,
//	    InterimResults: true,
//	})
//	if err != nil { return err }
//	defer stream.Close()
//
// See [ExampleClient_Connect] for the full runnable form.
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
	if opts != nil && len(opts.AdditionalQueryParams) > 0 {
		parsed, parseErr := url.Parse(dialURL)
		if parseErr != nil {
			return nil, fmt.Errorf("listen websocket: parse dial URL: %w", parseErr)
		}
		q := parsed.Query()
		for k, vs := range opts.AdditionalQueryParams {
			q.Del(k)
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		parsed.RawQuery = q.Encode()
		dialURL = parsed.String()
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
	stream := &Stream{transport: transport}
	stream.applyConfigSendDefaults(c.config)
	return stream, nil
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
//
// Calling Close cancels any outstanding SendAudioContext call and
// terminates the underlying WebSocket — there are no internal
// goroutines left running once Close returns successfully.
type Stream struct {
	transport    wstransport.Stream[spectypes.ClientStream, spectypes.ServerStream]
	maxFrameSize int
	sendTimeout  time.Duration
}

// SendAudio transmits a raw audio chunk as a binary WebSocket frame.
// The chunk's encoding must match the Encoding/SampleRate/Channels
// options used at Connect time. An empty chunk is treated by the
// server as equivalent to CloseStream.
//
// Returns ErrFrameTooLarge when the chunk exceeds Config.MaxFrameSizeBytes
// (if Client.WithConfig was used) and ErrSendTimeout when
// Config.SendTimeout elapses before the underlying write completes.
//
// Example:
//
//	buf := make([]byte, 4096)
//	for {
//	    n, err := mic.Read(buf)
//	    if err != nil { break }
//	    if err := stream.SendAudio(buf[:n]); err != nil { return err }
//	}
//
// See [ExampleStream_SendAudio] for the full runnable form.
func (s *Stream) SendAudio(data []byte) error {
	if s.maxFrameSize > 0 && len(data) > s.maxFrameSize {
		return fmt.Errorf("%w: %d > %d", ErrFrameTooLarge, len(data), s.maxFrameSize)
	}
	if s.sendTimeout <= 0 {
		return s.transport.Send(&spectypes.ClientStreamMemberAudio{
			Value: spectypes.AudioFrame{Data: data},
		})
	}
	done := make(chan error, 1)
	go func() {
		done <- s.transport.Send(&spectypes.ClientStreamMemberAudio{
			Value: spectypes.AudioFrame{Data: data},
		})
	}()
	timer := time.NewTimer(s.sendTimeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return ErrSendTimeout
	case err := <-done:
		return err
	}
}

// CloseStream sends a graceful end-of-audio marker. The server
// flushes pending results and emits a final MetadataEvent before
// closing the WebSocket.
//
// Example:
//
//	if err := stream.CloseStream(); err != nil { return err }
//	// keep calling stream.Recv until io.EOF for the final MetadataEvent
//
// See [ExampleStream_CloseStream] for the full runnable form.
func (s *Stream) CloseStream() error {
	return s.transport.Send(&spectypes.ClientStreamMemberCloseStream{
		Value: spectypes.CloseStream{},
	})
}

// Finalize forces the server to emit a final transcript for the
// open utterance. Pass channel = -1 to finalize all channels.
//
// Example:
//
//	if err := stream.Finalize(-1); err != nil { return err }
//
// See [ExampleStream_Finalize] for the full runnable form.
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
//
// Example:
//
//	ticker := time.NewTicker(5 * time.Second)
//	defer ticker.Stop()
//	for range ticker.C {
//	    if err := stream.KeepAlive(); err != nil { return err }
//	}
//
// See [ExampleStream_KeepAlive] for the full runnable form.
func (s *Stream) KeepAlive() error {
	return s.transport.Send(&spectypes.ClientStreamMemberKeepAlive{
		Value: spectypes.KeepAlive{},
	})
}

// Sync sends a client-side sync marker with the supplied id. The
// server echoes it back as a SyncEvent. Used to align test
// pipelines.
//
// Example:
//
//	if err := stream.Sync(42); err != nil { return err }
//	// next stream.Recv() will eventually yield *ws.SyncEvent{ID: 42}
//
// See [ExampleStream_Sync] for the full runnable form.
func (s *Stream) Sync(id int64) error {
	return s.transport.Send(&spectypes.ClientStreamMemberSync{
		Value: spectypes.ClientSync{Id: &id},
	})
}

// Recv blocks until the next event arrives or the connection
// closes. Returns nil for messages the SDK does not recognize so
// callers can continue the loop on forward-compat variants. Returns
// io.EOF (or other transport error) when the connection ends.
//
// Example:
//
//	for {
//	    event, err := stream.Recv()
//	    if err != nil { return err }
//	    switch e := event.(type) {
//	    case *ws.ResultsEvent:
//	        if e.IsFinal {
//	            fmt.Println(e.Channel.Alternatives[0].Transcript)
//	        }
//	    case *ws.MetadataEvent:
//	        return nil // session ended
//	    case *ws.ErrorEvent:
//	        return fmt.Errorf("stream error: %s", e.Description)
//	    }
//	}
//
// See [ExampleStream_Recv] for the full runnable form.
func (s *Stream) Recv() (Event, error) {
	msg, err := s.transport.Recv()
	if err != nil {
		return nil, err
	}
	return fromServerStream(msg), nil
}

// Close terminates the WebSocket connection. Idempotent; safe to
// call after CloseStream.
//
// Example:
//
//	stream, err := client.Connect(ctx, opts)
//	if err != nil { return err }
//	defer stream.Close()
//
// See [ExampleStream_Close] for the full runnable form.
func (s *Stream) Close() error {
	return s.transport.Close()
}
