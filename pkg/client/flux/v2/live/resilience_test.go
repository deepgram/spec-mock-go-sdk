// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Unit test for the generalized live-stream resilience (frame-size guard),
// ported from Listen-live. No network: a no-op transport + a Stream with
// maxFrameSize set directly.
package livev2

import (
	"errors"
	"io"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

type nopWire struct{}

func (nopWire) Send(spectypes.FluxClientStream) error     { return nil }
func (nopWire) Recv() (spectypes.FluxServerStream, error) { return nil, io.EOF }
func (nopWire) Close() error                              { return nil }

func TestSendAudio_FrameSizeGuard(t *testing.T) {
	s := &Stream{transport: nopWire{}, maxFrameSize: 10}
	if err := s.SendAudio(make([]byte, 11)); !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("oversized frame: got %v, want ErrFrameTooLarge", err)
	}
	if err := s.SendAudio(make([]byte, 5)); err != nil {
		t.Fatalf("in-limit frame: got %v, want nil", err)
	}
}

func TestApplyConfigSendDefaults(t *testing.T) {
	s := &Stream{transport: nopWire{}}
	s.applyConfigSendDefaults(&Config{MaxFrameSizeBytes: 42})
	if s.maxFrameSize != 42 {
		t.Fatalf("maxFrameSize = %d, want 42", s.maxFrameSize)
	}
}
