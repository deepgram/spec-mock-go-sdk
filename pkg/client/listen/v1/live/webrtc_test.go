// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// WebRTC transport: the contract (signaling endpoint, ICE servers, data
// channel) is modeled in @supportsTransports webrtc and generated as a typed
// config; the signaling handshake is an honest stub until Deepgram publishes a
// WebRTC endpoint. Verifies the spec-modeled defaults, transport selection, and
// the two documented stub errors.
package livev1

import (
	"context"
	"strings"
	"testing"
)

func TestWebRTCDefaultsAreModeled(t *testing.T) {
	cfg := DefaultWebRTCTransportConfig()
	if cfg.SignalingProtocol != "deepgram-webrtc-v1" {
		t.Fatalf("SignalingProtocol = %q, want deepgram-webrtc-v1", cfg.SignalingProtocol)
	}
	if cfg.SignalingURL != "" {
		t.Fatalf("SignalingURL default should be empty (caller-supplied), got %q", cfg.SignalingURL)
	}
	if len(cfg.ICEServers) != 1 || len(cfg.ICEServers[0].URLs) != 1 ||
		cfg.ICEServers[0].URLs[0] != "stun:stun.l.google.com:19302" {
		t.Fatalf("ICEServers = %+v", cfg.ICEServers)
	}
	if cfg.DataChannel.Label != "deepgram-listen" || !cfg.DataChannel.Ordered {
		t.Fatalf("DataChannel = %+v", cfg.DataChannel)
	}

	c := New(WithWebRTCTransport(cfg))
	if _, ok := c.transport.(webRTCBinding); !ok {
		t.Fatalf("WithWebRTCTransport should select webRTCBinding, got %T", c.transport)
	}
}

func TestWebRTCConnectRequiresSignalingURL(t *testing.T) {
	c := New(WithWebRTCTransport(DefaultWebRTCTransportConfig()))
	_, err := c.Connect(context.Background(), &LiveTranscriptionOptions{Model: "nova-3"})
	if err == nil || !strings.Contains(err.Error(), "SignalingURL is required") {
		t.Fatalf("expected SignalingURL-required error, got %v", err)
	}
}

func TestWebRTCConnectStubsHandshake(t *testing.T) {
	cfg := DefaultWebRTCTransportConfig()
	cfg.SignalingURL = "wss://example.invalid/webrtc"
	c := New(WithWebRTCTransport(cfg))
	_, err := c.Connect(context.Background(), &LiveTranscriptionOptions{Model: "nova-3"})
	if err == nil || !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("expected not-implemented handshake error, got %v", err)
	}
}
