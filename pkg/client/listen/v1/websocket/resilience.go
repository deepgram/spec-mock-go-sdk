// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import (
	"context"
	"errors"
	"fmt"
	"time"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

// Config carries facade-level resilience knobs that are NOT part of the
// Smithy wire surface. They tune how this SDK talks to the live
// transcription endpoint without changing what's sent on the wire.
//
// Zero-value Config is valid: every field falls back to a safe default
// documented per-field.
type Config struct {
	// MaxFrameSizeBytes caps the size of a single SendAudio chunk.
	// Audio chunks larger than this are rejected with ErrFrameTooLarge
	// before they reach the network. Zero (default) disables the
	// per-chunk limit; the underlying WebSocket library still enforces
	// its own per-frame maximum (1 MiB by default for dvonthenen/websocket).
	//
	// Recommended: 64 KiB for typical mic capture at 16-bit / 16 kHz
	// (~2 seconds of audio per frame). Larger frames increase end-to-end
	// latency.
	MaxFrameSizeBytes int

	// SendTimeout bounds how long Stream.SendAudio will wait for the
	// underlying write to complete. Zero (default) blocks indefinitely,
	// matching the historical behaviour. A non-zero value lets a stuck
	// network surface as ErrSendTimeout rather than a hung goroutine.
	//
	// Recommended: 5 * time.Second for interactive use; longer for
	// batch / pre-recorded scenarios where the caller has its own
	// supervision.
	SendTimeout time.Duration

	// ReconnectPolicy configures whether and how Connect-returned
	// Streams should attempt to recover from transport-level
	// disconnects. Nil (default) disables reconnection; the customer
	// receives the underlying transport error from Stream.Recv as
	// before.
	//
	// See ReconnectPolicy for the per-policy semantics.
	ReconnectPolicy *ReconnectPolicy
}

// ReconnectPolicy controls Stream-level auto-reconnect behaviour.
//
// When set on the Client's Config, the Stream returned by Connect will
// transparently re-dial the WebSocket on transport failure up to
// MaxAttempts times, with exponential backoff bounded by MaxBackoff.
// Customer code receives a *ReconnectEvent on Stream.Recv when a
// reconnect succeeds so it can resync state (e.g. replay tail audio).
//
// IMPORTANT: reconnect is best-effort. Session-level state on the
// Deepgram side is NOT preserved across reconnects — interim transcripts
// in flight at disconnect time are lost. Customers MUST be prepared to
// re-send the last N seconds of audio if continuity matters.
type ReconnectPolicy struct {
	// MaxAttempts caps the total number of re-dial attempts after a
	// single transport failure. Zero is treated as 3.
	MaxAttempts int

	// InitialBackoff is the delay before the first re-dial attempt.
	// Doubled on every subsequent attempt up to MaxBackoff. Zero is
	// treated as 500 * time.Millisecond.
	InitialBackoff time.Duration

	// MaxBackoff caps the per-attempt backoff. Zero is treated as
	// 30 * time.Second.
	MaxBackoff time.Duration
}

// ReconnectEvent is delivered via Stream.Recv when an automatic
// reconnect succeeds. Customer code typically responds by re-sending a
// short prefix of audio so the new session's transcript dovetails with
// the previous session's tail.
//
// This is a facade-only event type; it is NOT part of the Smithy wire
// surface. Server-emitted events continue to come through the standard
// Event variants in events.go.
type ReconnectEvent struct {
	Attempts        int
	BackoffApplied  time.Duration
	UnderlyingError error
}

func (*ReconnectEvent) isEvent() {}

// ErrFrameTooLarge is returned by Stream.SendAudio (and SendAudioContext)
// when the supplied chunk exceeds Config.MaxFrameSizeBytes.
var ErrFrameTooLarge = errors.New("listen websocket: audio frame exceeds Config.MaxFrameSizeBytes")

// ErrSendTimeout is returned by Stream.SendAudio (and SendAudioContext)
// when Config.SendTimeout elapses before the underlying write completes.
var ErrSendTimeout = errors.New("listen websocket: SendAudio exceeded Config.SendTimeout")

// WithConfig returns a copy of the Client with the supplied Config.
// The Config is consulted on Connect; subsequent edits do not affect
// already-open Streams.
//
// Example:
//
//	client := ws.NewWithDefaults().WithConfig(&ws.Config{
//	    MaxFrameSizeBytes: 64 * 1024,
//	    SendTimeout:       5 * time.Second,
//	    ReconnectPolicy: &ws.ReconnectPolicy{
//	        MaxAttempts: 5,
//	    },
//	})
//
// See [ExampleClient_WithConfig] for the full runnable form.
func (c *Client) WithConfig(cfg *Config) *Client {
	out := *c
	out.config = cfg
	return &out
}

// SendAudioContext is the context-aware sibling of SendAudio. The
// context bounds how long the call will wait for the underlying write;
// it is honoured AFTER any frame-size validation but BEFORE the
// transport's own write. Cancelling the context after the bytes have
// been handed to the transport does not recall them.
//
// Customers who do not need cancellation should keep using
// Stream.SendAudio.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//	if err := stream.SendAudioContext(ctx, chunk); err != nil {
//	    return err
//	}
func (s *Stream) SendAudioContext(ctx context.Context, data []byte) error {
	if s.maxFrameSize > 0 && len(data) > s.maxFrameSize {
		return fmt.Errorf("%w: %d > %d", ErrFrameTooLarge, len(data), s.maxFrameSize)
	}
	if ctx == nil {
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
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// applyConfigSendDefaults wraps SendAudio with the Config-derived
// behaviour (frame size, timeout). Used by Connect to set up the
// Stream's send pipeline. Not exported; callers go through SendAudio /
// SendAudioContext.
func (s *Stream) applyConfigSendDefaults(cfg *Config) {
	if cfg == nil {
		return
	}
	s.maxFrameSize = cfg.MaxFrameSizeBytes
	s.sendTimeout = cfg.SendTimeout
}

// effectiveReconnectPolicy returns the policy fields with defaults
// applied, suitable for the reconnect loop.
func (p *ReconnectPolicy) effective() (maxAttempts int, initial, max time.Duration) {
	maxAttempts = p.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = 3
	}
	initial = p.InitialBackoff
	if initial == 0 {
		initial = 500 * time.Millisecond
	}
	max = p.MaxBackoff
	if max == 0 {
		max = 30 * time.Second
	}
	return
}
