// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package deepgramtest

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	ws "github.com/deepgram/spec-mock-go-sdk/api/transport/websocket"
)

func Test_WS_MarshalClientStream_Audio_BinaryFrame(t *testing.T) {
	audio := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	msg := &spectypes.ClientStreamMemberAudio{
		Value: spectypes.AudioFrame{Data: audio},
	}
	payload, isBinary, err := ws.MarshalClientStream(msg)
	if err != nil {
		t.Fatalf("MarshalClientStream: %v", err)
	}
	if !isBinary {
		t.Fatal("audio variant should marshal as binary frame")
	}
	if !bytes.Equal(payload, audio) {
		t.Fatalf("payload = %v, want %v", payload, audio)
	}
}

func Test_WS_MarshalClientStream_CloseStream_TaggedJSON(t *testing.T) {
	msg := &spectypes.ClientStreamMemberCloseStream{
		Value: spectypes.CloseStream{},
	}
	payload, isBinary, err := ws.MarshalClientStream(msg)
	if err != nil {
		t.Fatalf("MarshalClientStream: %v", err)
	}
	if isBinary {
		t.Fatal("closeStream variant should marshal as JSON text frame, not binary")
	}
	var wire map[string]any
	if err := json.Unmarshal(payload, &wire); err != nil {
		t.Fatalf("payload not valid JSON: %v", err)
	}
	if wire["type"] != "CloseStream" {
		t.Fatalf(`wire["type"] = %v, want "CloseStream"`, wire["type"])
	}
}

func Test_WS_MarshalClientStream_KeepAlive_TaggedJSON(t *testing.T) {
	payload, isBinary, err := ws.MarshalClientStream(&spectypes.ClientStreamMemberKeepAlive{
		Value: spectypes.KeepAlive{},
	})
	if err != nil {
		t.Fatalf("MarshalClientStream: %v", err)
	}
	if isBinary {
		t.Fatal("keepAlive should be JSON text frame")
	}
	if !strings.Contains(string(payload), `"type":"KeepAlive"`) {
		t.Fatalf("payload = %s, want type discriminator KeepAlive", payload)
	}
}

func Test_WS_MarshalClientStream_Finalize_CarriesChannel(t *testing.T) {
	ch := int32(0)
	msg := &spectypes.ClientStreamMemberFinalize{
		Value: spectypes.Finalize{Channel: &ch},
	}
	payload, _, err := ws.MarshalClientStream(msg)
	if err != nil {
		t.Fatalf("MarshalClientStream: %v", err)
	}
	var wire map[string]any
	if err := json.Unmarshal(payload, &wire); err != nil {
		t.Fatalf("payload not valid JSON: %v", err)
	}
	if wire["type"] != "Finalize" {
		t.Fatalf(`type = %v, want "Finalize"`, wire["type"])
	}
	if c, ok := wire["channel"].(float64); !ok || int32(c) != ch {
		t.Fatalf("channel = %v, want %d", wire["channel"], ch)
	}
}

func Test_WS_MarshalClientStream_NilMessage_Errors(t *testing.T) {
	_, _, err := ws.MarshalClientStream(nil)
	if err == nil {
		t.Fatal("expected error on nil message, got nil")
	}
}

func Test_WS_UnmarshalServerStream_Results(t *testing.T) {
	wire := []byte(`{
		"type": "Results",
		"channel_index": [0, 1],
		"duration": 1.5,
		"start": 0.0,
		"is_final": true,
		"channel": {"alternatives": [{"transcript": "hello", "confidence": 0.95, "words": [], "languages": []}]},
		"metadata": {},
		"from_finalize": false
	}`)
	got, err := ws.UnmarshalServerStream(wire)
	if err != nil {
		t.Fatalf("UnmarshalServerStream: %v", err)
	}
	results, ok := got.(*spectypes.ServerStreamMemberResults)
	if !ok {
		t.Fatalf("got %T, want *ServerStreamMemberResults", got)
	}
	if results.Value.Channel.Alternatives[0].Transcript == nil ||
		*results.Value.Channel.Alternatives[0].Transcript != "hello" {
		t.Fatalf("transcript = %v, want 'hello'", results.Value.Channel.Alternatives[0].Transcript)
	}
}

func Test_WS_UnmarshalServerStream_UtteranceEnd(t *testing.T) {
	wire := []byte(`{
		"type": "UtteranceEnd",
		"channel": [0, 1],
		"last_word_end": 2.5
	}`)
	got, err := ws.UnmarshalServerStream(wire)
	if err != nil {
		t.Fatalf("UnmarshalServerStream: %v", err)
	}
	ue, ok := got.(*spectypes.ServerStreamMemberUtteranceEnd)
	if !ok {
		t.Fatalf("got %T, want *ServerStreamMemberUtteranceEnd", got)
	}
	if ue.Value.LastWordEnd == nil || *ue.Value.LastWordEnd != 2.5 {
		t.Fatalf("LastWordEnd = %v, want 2.5", ue.Value.LastWordEnd)
	}
}

func Test_WS_UnmarshalServerStream_UnknownType_Errors(t *testing.T) {
	_, err := ws.UnmarshalServerStream([]byte(`{"type":"NotAThing"}`))
	if err == nil {
		t.Fatal("expected error on unknown type, got nil")
	}
}

func Test_WS_UnmarshalClientStream_BinaryFrame_RoutesToAudio(t *testing.T) {
	audio := []byte{0xAA, 0xBB, 0xCC}
	got, err := ws.UnmarshalClientStream(audio, true)
	if err != nil {
		t.Fatalf("UnmarshalClientStream: %v", err)
	}
	audioVar, ok := got.(*spectypes.ClientStreamMemberAudio)
	if !ok {
		t.Fatalf("got %T, want *ClientStreamMemberAudio", got)
	}
	if !bytes.Equal(audioVar.Value.Data, audio) {
		t.Fatalf("audio data = %v, want %v", audioVar.Value.Data, audio)
	}
}

func Test_WS_CloseCodeConstants(t *testing.T) {
	if ws.CloseCodeClientTimeout != 1011 {
		t.Fatalf("CloseCodeClientTimeout = %d, want 1011", ws.CloseCodeClientTimeout)
	}
	if ws.CloseCodeDriverTimeout != 1011 {
		t.Fatalf("CloseCodeDriverTimeout = %d, want 1011", ws.CloseCodeDriverTimeout)
	}
	if ws.CloseCodeCodec != 1008 {
		t.Fatalf("CloseCodeCodec = %d, want 1008", ws.CloseCodeCodec)
	}
}
