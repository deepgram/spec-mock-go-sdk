// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Converters: spectypes.* (generated, pointer-typed, in api/types) ->
// idiomatic value-typed structs in this package. Routers in
// pkg/api/listen/v1/websocket/{chan,callback}_router.go call these at
// the fan-out boundary so customer code receives value-typed messages
// without ever nil-checking pointer fields.

package interfacesv1

import (
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefF32(p *float32) float64 {
	if p == nil {
		return 0
	}
	return float64(*p)
}

func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func derefInt32(p *int32) int {
	if p == nil {
		return 0
	}
	return int(*p)
}

func derefInt64ToPtrInt(p *int64) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}

// FromStreamingResponse converts the generated Results variant payload
// into the idiomatic MessageResponse customers see.
func FromStreamingResponse(in *spectypes.StreamingResponse) *MessageResponse {
	if in == nil {
		return nil
	}
	out := &MessageResponse{
		Duration:     derefF32(in.Duration),
		IsFinal:      derefBool(in.IsFinal),
		FromFinalize: derefBool(in.FromFinalize),
		SpeechFinal:  derefBool(in.SpeechFinal),
		Start:        derefF32(in.Start),
		Type:         "Results",
	}
	if len(in.ChannelIndex) > 0 {
		out.ChannelIndex = make([]int, len(in.ChannelIndex))
		for i, v := range in.ChannelIndex {
			out.ChannelIndex[i] = int(v)
		}
	}
	if len(in.Channel.Alternatives) > 0 {
		out.Channel.Alternatives = make([]Alternative, len(in.Channel.Alternatives))
		for i, a := range in.Channel.Alternatives {
			out.Channel.Alternatives[i] = convertAlternative(&a)
		}
	}
	if in.Metadata != nil {
		out.Metadata = Metadata{
			RequestID: derefStr(in.Metadata.RequestId),
			ModelUUID: derefStr(in.Metadata.ModelUuid),
			Extra:     in.Metadata.Extra,
		}
		if in.Metadata.ModelInfo != nil {
			out.Metadata.ModelInfo = ModelInfo{
				Arch:    derefStr(in.Metadata.ModelInfo.Arch),
				Name:    derefStr(in.Metadata.ModelInfo.Name),
				Version: derefStr(in.Metadata.ModelInfo.Version),
			}
		}
	}
	return out
}

func convertAlternative(in *spectypes.Alternative) Alternative {
	out := Alternative{
		Confidence: derefF32(in.Confidence),
		Transcript: derefStr(in.Transcript),
		Languages:  in.Languages,
	}
	if len(in.Words) > 0 {
		out.Words = make([]Word, len(in.Words))
		for i, w := range in.Words {
			out.Words[i] = Word{
				Confidence:     derefF32(w.Confidence),
				End:            derefF32(w.End),
				PunctuatedWord: derefStr(w.PunctuatedWord),
				Start:          derefF32(w.Start),
				Word:           derefStr(w.Word),
				Speaker:        derefInt64ToPtrInt(w.Speaker),
				Language:       derefStr(w.Language),
			}
		}
	}
	return out
}

// FromWsMetadata converts the generated Metadata variant payload into
// the idiomatic MetadataResponse customers see.
func FromWsMetadata(in *spectypes.WsMetadata) *MetadataResponse {
	if in == nil {
		return nil
	}
	out := &MetadataResponse{
		Channels:       derefInt32(in.Channels),
		Created:        derefStr(in.Created),
		Duration:       derefF32(in.Duration),
		RequestID:      derefStr(in.RequestId),
		Sha256:         derefStr(in.Sha256),
		TransactionKey: derefStr(in.TransactionKey),
		Type:           "Metadata",
		Extra:          in.Extra,
	}
	if len(in.ModelInfo) > 0 {
		out.ModelInfo = make(map[string]ModelInfo, len(in.ModelInfo))
		for k, v := range in.ModelInfo {
			out.ModelInfo[k] = ModelInfo{
				Arch:    derefStr(v.Arch),
				Name:    derefStr(v.Name),
				Version: derefStr(v.Version),
			}
		}
	}
	return out
}

// FromUtteranceEnd converts the generated UtteranceEnd variant payload
// into the idiomatic UtteranceEndResponse customers see.
func FromUtteranceEnd(in *spectypes.UtteranceEnd) *UtteranceEndResponse {
	if in == nil {
		return nil
	}
	out := &UtteranceEndResponse{
		Type:        "UtteranceEnd",
		LastWordEnd: derefF32(in.LastWordEnd),
	}
	if len(in.Channel) > 0 {
		out.Channel = make([]int, len(in.Channel))
		for i, v := range in.Channel {
			out.Channel[i] = int(v)
		}
	}
	return out
}

// FromSpeechStarted converts the generated SpeechStarted variant
// payload into the idiomatic SpeechStartedResponse customers see.
func FromSpeechStarted(in *spectypes.SpeechStarted) *SpeechStartedResponse {
	if in == nil {
		return nil
	}
	out := &SpeechStartedResponse{
		Type:      "SpeechStarted",
		Timestamp: derefF32(in.Timestamp),
	}
	if len(in.Channel) > 0 {
		out.Channel = make([]int, len(in.Channel))
		for i, v := range in.Channel {
			out.Channel[i] = int(v)
		}
	}
	return out
}
