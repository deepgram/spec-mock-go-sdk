// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package livev1

import (
	"context"
	"errors"
	"testing"
	"time"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

type mockWireStream struct {
	sent        []spectypes.ClientStream
	sendStarted chan struct{}
	releaseSend chan struct{}
}

func (m *mockWireStream) Send(msg spectypes.ClientStream) error {
	if m.sendStarted != nil {
		close(m.sendStarted)
	}
	if m.releaseSend != nil {
		<-m.releaseSend
	}
	m.sent = append(m.sent, msg)
	return nil
}

func (m *mockWireStream) Recv() (spectypes.ServerStream, error) { return nil, nil }
func (m *mockWireStream) Close() error                          { return nil }

func TestSendAudio_FrameSizeExceeded(t *testing.T) {
	stream := &Stream{transport: &mockWireStream{}, maxFrameSize: 4}
	err := stream.SendAudio([]byte("12345"))
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("SendAudio error = %v, want ErrFrameTooLarge", err)
	}
}

func TestSendAudio_RespectsContext(t *testing.T) {
	transport := &mockWireStream{sendStarted: make(chan struct{}), releaseSend: make(chan struct{})}
	stream := &Stream{transport: transport}
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() { errCh <- stream.SendAudioContext(ctx, []byte("audio")) }()
	<-transport.sendStarted
	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("SendAudioContext error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("SendAudioContext did not return after context cancellation")
	}
	close(transport.releaseSend)
}

func TestConfig_Defaults(t *testing.T) {
	stream := &Stream{}
	stream.applyConfigSendDefaults(nil)
	if stream.maxFrameSize != 0 {
		t.Fatalf("maxFrameSize = %d, want 0", stream.maxFrameSize)
	}
	if stream.sendTimeout != 0 {
		t.Fatalf("sendTimeout = %v, want 0", stream.sendTimeout)
	}
}

func TestReconnectPolicy_Defaults(t *testing.T) {
	policy := &ReconnectPolicy{}
	config := &Config{ReconnectPolicy: policy}
	client := New(WithConfig(config))
	if client.config.ReconnectPolicy != policy {
		t.Fatal("ReconnectPolicy was not retained on Config")
	}
	var _ Event = (*ReconnectEvent)(nil)
}
