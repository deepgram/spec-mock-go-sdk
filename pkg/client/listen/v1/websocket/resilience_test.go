// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import (
	"errors"
	"testing"
	"time"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

// Tests for the resilience surface. None of these touch the network;
// they exercise the SDK-side wiring of Config -> Stream behaviour.

type mockTransport struct {
	sent [][]byte
}

func (m *mockTransport) Send(msg spectypes.ClientStream) error {
	if audio, ok := msg.(*spectypes.ClientStreamMemberAudio); ok {
		m.sent = append(m.sent, audio.Value.Data)
	}
	return nil
}

func (m *mockTransport) Recv() (spectypes.ServerStream, error) {
	return nil, nil
}

func (m *mockTransport) Close() error { return nil }

func TestConfig_WithConfigStoresOnClient(t *testing.T) {
	cfg := &Config{
		MaxFrameSizeBytes: 64 * 1024,
		SendTimeout:       5 * time.Second,
	}
	client := New("test-key", "").WithConfig(cfg)
	if client.config != cfg {
		t.Fatalf("Client.WithConfig did not store config; got %+v", client.config)
	}
}

func TestConfig_WithConfigReturnsCopy(t *testing.T) {
	original := New("test-key", "")
	configured := original.WithConfig(&Config{MaxFrameSizeBytes: 1024})
	if original.config != nil {
		t.Fatalf("Client.WithConfig mutated receiver; original.config = %+v, expected nil", original.config)
	}
	if configured.config == nil {
		t.Fatal("Client.WithConfig returned a Client with nil config")
	}
}

func TestStream_SendAudioRejectsOversizedFrame(t *testing.T) {
	stream := &Stream{maxFrameSize: 100}
	err := stream.SendAudio(make([]byte, 200))
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("expected ErrFrameTooLarge, got %v", err)
	}
}

func TestStream_SendAudioAllowsExactSizeFrame(t *testing.T) {
	stream := &Stream{maxFrameSize: 100, transport: &mockTransport{}}
	err := stream.SendAudio(make([]byte, 100))
	if err != nil {
		t.Fatalf("exact-size frame should not error; got %v", err)
	}
}

func TestStream_SendAudioZeroLimitAllowsAnySize(t *testing.T) {
	stream := &Stream{maxFrameSize: 0, transport: &mockTransport{}}
	err := stream.SendAudio(make([]byte, 10_000_000))
	if err != nil {
		t.Fatalf("zero maxFrameSize should allow any size; got %v", err)
	}
}

func TestReconnectPolicy_DefaultsApplied(t *testing.T) {
	p := &ReconnectPolicy{}
	maxAttempts, initial, max := p.effective()
	if maxAttempts != 3 {
		t.Errorf("default MaxAttempts: want 3, got %d", maxAttempts)
	}
	if initial != 500*time.Millisecond {
		t.Errorf("default InitialBackoff: want 500ms, got %v", initial)
	}
	if max != 30*time.Second {
		t.Errorf("default MaxBackoff: want 30s, got %v", max)
	}
}

func TestReconnectPolicy_CustomValuesPreserved(t *testing.T) {
	p := &ReconnectPolicy{
		MaxAttempts:    5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
	}
	maxAttempts, initial, max := p.effective()
	if maxAttempts != 5 || initial != 100*time.Millisecond || max != 10*time.Second {
		t.Errorf("custom values not preserved: %d / %v / %v", maxAttempts, initial, max)
	}
}

func TestReconnectEvent_ImplementsEventInterface(t *testing.T) {
	var _ Event = (*ReconnectEvent)(nil)
}
