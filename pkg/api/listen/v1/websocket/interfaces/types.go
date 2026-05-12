// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Idiomatic value-typed shapes for the listen WS surface. Customers
// receive these from the channel- and callback-based handlers. Same
// facade pattern as pkg/client/listen/v1/rest/response.go: the
// generated spectypes.* shapes parse the wire bytes, then the routers
// convert to these value-typed structs (see convert.go in this
// package) before fanning out. Field renames, nullability changes, and
// pointer vs value typing in api/types are all absorbed here.

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces"
)

// Option types - configuration sent client-to-server, not part of the
// inbound wire surface.
type LiveOptions interfaces.LiveTranscriptionOptions
type LiveTranscriptionOptions interfaces.LiveTranscriptionOptions

// MessageType is the type-discriminator header. Kept for back-compat
// with code that peeks at the type field directly; new code should use
// ws.UnmarshalServerStream which handles the discriminator internally.
type MessageType struct {
	Type string `json:"type"`
}

// Word is a single word in a transcript.
type Word struct {
	Confidence     float64 `json:"confidence,omitempty"`
	End            float64 `json:"end,omitempty"`
	PunctuatedWord string  `json:"punctuated_word,omitempty"`
	Start          float64 `json:"start,omitempty"`
	Word           string  `json:"word,omitempty"`
	Speaker        *int    `json:"speaker,omitempty"`
	Language       string  `json:"language,omitempty"`
}

// Alternative is one ranked alternative interpretation.
type Alternative struct {
	Confidence float64  `json:"confidence,omitempty"`
	Transcript string   `json:"transcript,omitempty"`
	Words      []Word   `json:"words,omitempty"`
	Languages  []string `json:"languages,omitempty"`
}

// Channel is the per-channel transcript container.
type Channel struct {
	Alternatives []Alternative `json:"alternatives,omitempty"`
}

// ModelInfo describes the model that produced the transcript.
type ModelInfo struct {
	Arch    string `json:"arch,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Metadata is per-interim-result metadata embedded in MessageResponse.
type Metadata struct {
	Extra     map[string]string `json:"extra,omitempty"`
	ModelInfo ModelInfo         `json:"model_info,omitempty"`
	ModelUUID string            `json:"model_uuid,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// MessageResponse is the Results message variant - interim or final
// transcription result.
type MessageResponse struct {
	Channel      Channel  `json:"channel,omitempty"`
	ChannelIndex []int    `json:"channel_index,omitempty"`
	Duration     float64  `json:"duration,omitempty"`
	IsFinal      bool     `json:"is_final,omitempty"`
	FromFinalize bool     `json:"from_finalize,omitempty"`
	Metadata     Metadata `json:"metadata,omitempty"`
	SpeechFinal  bool     `json:"speech_final,omitempty"`
	Start        float64  `json:"start,omitempty"`
	Type         string   `json:"type,omitempty"`
}

// MetadataResponse is the per-session Metadata message variant
// (emitted at session end).
type MetadataResponse struct {
	Channels       int                  `json:"channels,omitempty"`
	Created        string               `json:"created,omitempty"`
	Duration       float64              `json:"duration,omitempty"`
	ModelInfo      map[string]ModelInfo `json:"model_info,omitempty"`
	Models         []string             `json:"models,omitempty"`
	RequestID      string               `json:"request_id,omitempty"`
	Sha256         string               `json:"sha256,omitempty"`
	TransactionKey string               `json:"transaction_key,omitempty"`
	Type           string               `json:"type,omitempty"`
	Extra          map[string]string    `json:"extra,omitempty"`
}

// UtteranceEndResponse signals a VAD-derived utterance boundary.
type UtteranceEndResponse struct {
	Type        string  `json:"type,omitempty"`
	Channel     []int   `json:"channel,omitempty"`
	LastWordEnd float64 `json:"last_word_end,omitempty"`
}

// SpeechStartedResponse signals speech detection after silence.
type SpeechStartedResponse struct {
	Type      string  `json:"type,omitempty"`
	Channel   []int   `json:"channel,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
}

// Connection-level events - not Smithy-modelled, left as before.
type OpenResponse = commoninterfaces.OpenResponse
type CloseResponse = commoninterfaces.CloseResponse

// SDK error type - not a ServerStream variant. Left as before.
type ErrorResponse = interfaces.DeepgramError
