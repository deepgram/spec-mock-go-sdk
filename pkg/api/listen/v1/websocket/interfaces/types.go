// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// This file used to declare hand-written wire types (Word / Alternative /
// Channel / Metadata / MessageResponse / MetadataResponse / etc.) that
// duplicated the JSON shapes already defined by deepgram/spec. After the
// listen WS production-faithful rewire, every wire-related type here is
// a Go type alias for its generated counterpart in api/types. Connection-
// level events (OpenResponse / CloseResponse) and the SDK's generic
// ErrorResponse stay as before because they aren't part of the Smithy-
// modelled wire surface.
//
// Field-level breaking changes vs the original hand-written types:
//   - all fields are now pointers (smithy-go's convention for nullable
//     wire members), e.g. `*string` rather than `string`
//   - floats are `*float32` rather than `float64` to match the Smithy
//     Float primitive precision
//
// Customer code accessing these fields needs to nil-check / dereference;
// the legacy default handlers and examples in this repo do.

package interfacesv1

import (
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	commoninterfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces"
)

// Option types - configuration carried client-to-server, not part of the
// wire message surface. Left as-is.
type LiveOptions interfaces.LiveTranscriptionOptions
type LiveTranscriptionOptions interfaces.LiveTranscriptionOptions

// MessageType is the discriminator-only header used by legacy code to peek
// at JSON before unmarshalling into a typed struct. Retained for
// back-compat with code that imports it directly; new code should use
// api/transport/websocket.UnmarshalServerStream which handles the
// discriminator automatically.
type MessageType struct {
	Type string `json:"type"`
}

// Wire shapes - aliased to the spec-generated equivalents in api/types.
// These four were the only message envelopes the router parses on the
// inbound side; everything else (Word / Alternative / Channel / etc.) is
// reachable as fields of these types.
type MessageResponse = spectypes.StreamingResponse
type MetadataResponse = spectypes.WsMetadata
type UtteranceEndResponse = spectypes.UtteranceEnd
type SpeechStartedResponse = spectypes.SpeechStarted

// Transcript-shape aliases for consumers that imported these names
// directly. All field types match api/types verbatim - they ARE api/types.
type Word = spectypes.Word
type Alternative = spectypes.Alternative
type Channel = spectypes.Channel
type ModelInfo = spectypes.ModelInfo
type Metadata = spectypes.WsMetadata

// Connection-level events - not Smithy-modelled, left as before.
type OpenResponse = commoninterfaces.OpenResponse
type CloseResponse = commoninterfaces.CloseResponse

// SDK error type - not a ServerStream variant. Left as before.
type ErrorResponse = interfaces.DeepgramError
