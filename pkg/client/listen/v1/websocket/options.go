// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import "net/url"

// LiveTranscriptionOptions are the customer-facing options for the
// /v1/listen WebSocket stream. Each exported field corresponds to a
// single @httpQuery member on the generated spectypes.StreamInput
// wire shape; the client exposes idiomatic Go value types so callers
// do not nil-check primitives.
//
// Zero values mean "not set" and leave the corresponding wire query
// parameter unset.
//
// The set of fields here is the public surface of /v1/listen
// streaming. Any @internal-tagged parameter on the wire input shape
// is intentionally absent. The SecWebSocketProtocol field on
// StreamInput is also absent: it carries browser-side credential
// delivery and is owned by the Client, not by per-call options.
type LiveTranscriptionOptions struct {
	Channels        int      `json:"channels,omitempty"          schema:"channels,omitempty"`

	// Deprecated: Legacy flag. Prefer DiarizeModel for explicit model selection.
	// Mutually exclusive with DiarizeModel.
	Diarize         bool     `json:"diarize,omitempty"           schema:"diarize,omitempty"`

	// Deprecated: Legacy diarization-model selector. Prefer DiarizeModel.
	DiarizeVersion  string   `json:"diarize_version,omitempty"   schema:"diarize_version,omitempty"`

	Encoding        string   `json:"encoding,omitempty"          schema:"encoding,omitempty"`
	Endpointing     int      `json:"endpointing,omitempty"       schema:"endpointing,omitempty"`
	InterimResults  bool     `json:"interim_results,omitempty"   schema:"interim_results,omitempty"`
	Keyterm         []string `json:"keyterm,omitempty"           schema:"keyterm,omitempty"`
	Keywords        []string `json:"keywords,omitempty"          schema:"keywords,omitempty"`
	Language        string   `json:"language,omitempty"          schema:"language,omitempty"`

	// Deprecated: Prefer MipOptOut. LogData is recognized for backward
	// compatibility.
	LogData         bool     `json:"log_data,omitempty"          schema:"log_data,omitempty"`

	MipOptOut       bool     `json:"mip_opt_out,omitempty"       schema:"mip_opt_out,omitempty"`
	Model           string   `json:"model,omitempty"             schema:"model,omitempty"`
	Multichannel    bool     `json:"multichannel,omitempty"      schema:"multichannel,omitempty"`
	ProfanityFilter bool     `json:"profanity_filter,omitempty"  schema:"profanity_filter,omitempty"`
	Punctuate       bool     `json:"punctuate,omitempty"         schema:"punctuate,omitempty"`
	Redact          []string `json:"redact,omitempty"            schema:"redact,omitempty"`
	SampleRate      int      `json:"sample_rate,omitempty"       schema:"sample_rate,omitempty"`
	Search          []string `json:"search,omitempty"            schema:"search,omitempty"`
	SmartFormat     bool     `json:"smart_format,omitempty"      schema:"smart_format,omitempty"`
	Tag             []string `json:"tag,omitempty"               schema:"tag,omitempty"`
	UtteranceEndMs  int      `json:"utterance_end_ms,omitempty"  schema:"utterance_end_ms,omitempty"`
	VadEvents       bool     `json:"vad_events,omitempty"        schema:"vad_events,omitempty"`

	// Deprecated: Use Endpointing instead. Rejected when Endpointing is
	// also set as an integer.
	VadTurnoff      int      `json:"vad_turnoff,omitempty"       schema:"vad_turnoff,omitempty"`

	Version         string   `json:"version,omitempty"           schema:"version,omitempty"`

	// AdditionalQueryParams carries arbitrary query parameters to send
	// on the WebSocket upgrade URL. Keys are query-parameter names,
	// values are raw string values. Multiple values per key produce
	// repeated ?key=v1&key=v2 entries. When a key collides with one of
	// the typed fields above, AdditionalQueryParams wins.
	//
	// Use this when Deepgram ships a new streaming parameter on the API
	// that the SDK has not yet been updated to expose with a typed
	// field.
	AdditionalQueryParams url.Values `json:"-" schema:"-"`
}
