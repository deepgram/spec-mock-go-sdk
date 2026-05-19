// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

import (
	"net/url"
	"strconv"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

// optionsToStreamInput translates an idiomatic *LiveTranscriptionOptions
// into a *spectypes.StreamInput suitable for URL-encoding into the
// WebSocket dial URL. Zero-valued facade fields leave the generated
// pointer as nil so the parameter is omitted from the query string.
//
// Every @httpQuery member of spectypes.StreamInput has a wiring block
// here. @internal-tagged stem parameters that the spec excludes from
// the public surface are intentionally not represented on
// LiveTranscriptionOptions and therefore have no wiring block.
func optionsToStreamInput(o *LiveTranscriptionOptions) *spectypes.StreamInput {
	in := &spectypes.StreamInput{}
	if o == nil {
		return in
	}
	if o.Channels != 0 {
		v := int32(o.Channels)
		in.Channels = &v
	}
	if o.Diarize {
		v := o.Diarize
		in.Diarize = &v
	}
	if o.DiarizeVersion != "" {
		v := o.DiarizeVersion
		in.DiarizeVersion = &v
	}
	if o.Encoding != "" {
		v := o.Encoding
		in.Encoding = &v
	}
	if o.Endpointing != 0 {
		v := int32(o.Endpointing)
		in.Endpointing = &v
	}
	if o.InterimResults {
		v := o.InterimResults
		in.InterimResults = &v
	}
	if len(o.Keyterm) > 0 {
		in.Keyterm = o.Keyterm
	}
	if len(o.Keywords) > 0 {
		in.Keywords = o.Keywords
	}
	if o.Language != "" {
		v := o.Language
		in.Language = &v
	}
	if o.LogData {
		v := o.LogData
		in.LogData = &v
	}
	if o.MipOptOut {
		v := o.MipOptOut
		in.MipOptOut = &v
	}
	if o.Model != "" {
		v := o.Model
		in.Model = &v
	}
	if o.Multichannel {
		v := o.Multichannel
		in.Multichannel = &v
	}
	if o.ProfanityFilter {
		v := o.ProfanityFilter
		in.ProfanityFilter = &v
	}
	if o.Punctuate {
		v := o.Punctuate
		in.Punctuate = &v
	}
	if len(o.Redact) > 0 {
		in.Redact = o.Redact
	}
	if o.SampleRate != 0 {
		v := int32(o.SampleRate)
		in.SampleRate = &v
	}
	if len(o.Search) > 0 {
		in.Search = o.Search
	}
	if o.SmartFormat {
		v := o.SmartFormat
		in.SmartFormat = &v
	}
	if len(o.Tag) > 0 {
		in.Tag = o.Tag
	}
	if o.UtteranceEndMs != 0 {
		v := int32(o.UtteranceEndMs)
		in.UtteranceEndMs = &v
	}
	if o.VadEvents {
		v := o.VadEvents
		in.VadEvents = &v
	}
	if o.VadTurnoff != 0 {
		v := int32(o.VadTurnoff)
		in.VadTurnoff = &v
	}
	if o.Version != "" {
		v := o.Version
		in.Version = &v
	}
	return in
}

// streamInputQueryString encodes a *spectypes.StreamInput into a
// URL-encoded query string. Mirrors the wire names from the
// listen-live @httpQuery annotations one-for-one. Used by Client.Connect
// to assemble the WebSocket dial URL.
//
// This is the streaming equivalent of api/transport/http.Invoke's
// reflection-based query building. Streaming doesn't get the generated
// route metadata (no @http binding on the streaming operation), so the
// query string is built explicitly here.
func streamInputQueryString(in *spectypes.StreamInput) string {
	q := url.Values{}
	if in.Channels != nil {
		q.Set("channels", strconv.FormatInt(int64(*in.Channels), 10))
	}
	if in.Diarize != nil {
		q.Set("diarize", strconv.FormatBool(*in.Diarize))
	}
	if in.DiarizeVersion != nil {
		q.Set("diarize_version", *in.DiarizeVersion)
	}
	if in.Encoding != nil {
		q.Set("encoding", *in.Encoding)
	}
	if in.Endpointing != nil {
		q.Set("endpointing", strconv.FormatInt(int64(*in.Endpointing), 10))
	}
	if in.InterimResults != nil {
		q.Set("interim_results", strconv.FormatBool(*in.InterimResults))
	}
	for _, v := range in.Keyterm {
		q.Add("keyterm", v)
	}
	for _, v := range in.Keywords {
		q.Add("keywords", v)
	}
	if in.Language != nil {
		q.Set("language", *in.Language)
	}
	if in.LogData != nil {
		q.Set("log_data", strconv.FormatBool(*in.LogData))
	}
	if in.MipOptOut != nil {
		q.Set("mip_opt_out", strconv.FormatBool(*in.MipOptOut))
	}
	if in.Model != nil {
		q.Set("model", *in.Model)
	}
	if in.Multichannel != nil {
		q.Set("multichannel", strconv.FormatBool(*in.Multichannel))
	}
	if in.ProfanityFilter != nil {
		q.Set("profanity_filter", strconv.FormatBool(*in.ProfanityFilter))
	}
	if in.Punctuate != nil {
		q.Set("punctuate", strconv.FormatBool(*in.Punctuate))
	}
	for _, v := range in.Redact {
		q.Add("redact", v)
	}
	if in.SampleRate != nil {
		q.Set("sample_rate", strconv.FormatInt(int64(*in.SampleRate), 10))
	}
	for _, v := range in.Search {
		q.Add("search", v)
	}
	if in.SmartFormat != nil {
		q.Set("smart_format", strconv.FormatBool(*in.SmartFormat))
	}
	for _, v := range in.Tag {
		q.Add("tag", v)
	}
	if in.UtteranceEndMs != nil {
		q.Set("utterance_end_ms", strconv.FormatInt(int64(*in.UtteranceEndMs), 10))
	}
	if in.VadEvents != nil {
		q.Set("vad_events", strconv.FormatBool(*in.VadEvents))
	}
	if in.VadTurnoff != nil {
		q.Set("vad_turnoff", strconv.FormatInt(int64(*in.VadTurnoff), 10))
	}
	if in.Version != nil {
		q.Set("version", *in.Version)
	}
	return q.Encode()
}

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

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func int16SliceToInt(in []int16) []int {
	if len(in) == 0 {
		return nil
	}
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}

func int64ToPtrInt(p *int64) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}

func f32ToPtrF64(p *float32) *float64 {
	if p == nil {
		return nil
	}
	v := float64(*p)
	return &v
}

// fromServerStream converts a generated ServerStream union member into
// a customer-facing Event value. Returns nil for unknown variants;
// callers should treat nil as "ignored / unrecognized message" and
// continue the recv loop.
func fromServerStream(msg spectypes.ServerStream) Event {
	switch m := msg.(type) {
	case *spectypes.ServerStreamMemberResults:
		return fromStreamingResponse(&m.Value)
	case *spectypes.ServerStreamMemberMetadata:
		return fromWsMetadata(&m.Value)
	case *spectypes.ServerStreamMemberSpeechStarted:
		return fromSpeechStarted(&m.Value)
	case *spectypes.ServerStreamMemberUtteranceEnd:
		return fromUtteranceEnd(&m.Value)
	case *spectypes.ServerStreamMemberError:
		return fromWsError(&m.Value)
	case *spectypes.ServerStreamMemberSync:
		return fromServerSync(&m.Value)
	default:
		return nil
	}
}

func fromStreamingResponse(in *spectypes.StreamingResponse) *ResultsEvent {
	out := &ResultsEvent{
		ChannelIndex: int16SliceToInt(in.ChannelIndex),
		Duration:     derefF32(in.Duration),
		FromFinalize: derefBool(in.FromFinalize),
		IsFinal:      derefBool(in.IsFinal),
		Start:        derefF32(in.Start),
		SpeechFinal:  derefBool(in.SpeechFinal),
		Tag:          derefStr(in.Tag),
	}
	if in.Channel != nil {
		out.Channel = fromChannel(in.Channel)
	}
	if in.Metadata != nil {
		out.Metadata = fromInterimMetadata(in.Metadata)
	}
	if len(in.Entities) > 0 {
		out.Entities = make([]Entity, len(in.Entities))
		for i, e := range in.Entities {
			out.Entities[i] = Entity{
				Label:      derefStr(e.Label),
				Value:      derefStr(e.Value),
				Confidence: derefF32(e.Confidence),
				StartWord:  derefInt32(e.StartWord),
				EndWord:    derefInt32(e.EndWord),
			}
		}
	}
	return out
}

func fromChannel(in *spectypes.Channel) *Channel {
	out := &Channel{
		DetectedLanguage:   derefStr(in.DetectedLanguage),
		LanguageConfidence: derefF32(in.LanguageConfidence),
	}
	if len(in.Search) > 0 {
		out.Search = make([]Search, len(in.Search))
		for i, s := range in.Search {
			out.Search[i] = Search{Query: derefStr(s.Query)}
			if len(s.Hits) > 0 {
				out.Search[i].Hits = make([]Hit, len(s.Hits))
				for j, h := range s.Hits {
					out.Search[i].Hits[j] = Hit{
						Confidence: derefF32(h.Confidence),
						Start:      derefF32(h.Start),
						End:        derefF32(h.End),
						Snippet:    derefStr(h.Snippet),
					}
				}
			}
		}
	}
	if len(in.Alternatives) > 0 {
		out.Alternatives = make([]Alternative, len(in.Alternatives))
		for i, a := range in.Alternatives {
			out.Alternatives[i] = Alternative{
				Transcript: derefStr(a.Transcript),
				Confidence: derefF32(a.Confidence),
				Languages:  a.Languages,
			}
			if len(a.Words) > 0 {
				out.Alternatives[i].Words = make([]Word, len(a.Words))
				for j, w := range a.Words {
					out.Alternatives[i].Words[j] = Word{
						Word:              derefStr(w.Word),
						Start:             derefF32(w.Start),
						End:               derefF32(w.End),
						Confidence:        derefF32(w.Confidence),
						Speaker:           int64ToPtrInt(w.Speaker),
						SpeakerConfidence: f32ToPtrF64(w.SpeakerConfidence),
						PunctuatedWord:    derefStr(w.PunctuatedWord),
						Language:          derefStr(w.Language),
					}
				}
			}
		}
	}
	return out
}

func fromInterimMetadata(in *spectypes.InterimMetadata) *InterimMetadata {
	out := &InterimMetadata{
		RequestID: derefStr(in.RequestId),
		ModelUUID: derefStr(in.ModelUuid),
		Extra:     in.Extra,
	}
	if in.ModelInfo != nil {
		out.ModelInfo = &ModelInfo{
			Name:    derefStr(in.ModelInfo.Name),
			Version: derefStr(in.ModelInfo.Version),
			Arch:    derefStr(in.ModelInfo.Arch),
		}
	}
	return out
}

func fromWsMetadata(in *spectypes.WsMetadata) *MetadataEvent {
	out := &MetadataEvent{
		RequestID: derefStr(in.RequestId),
		Channels:  derefInt32(in.Channels),
		Created:   derefStr(in.Created),
		Duration:  derefF32(in.Duration),
		Sha256:    derefStr(in.Sha256),
		Extra:     in.Extra,
	}
	if len(in.ModelInfo) > 0 {
		out.ModelInfo = make(map[string]ModelInfo, len(in.ModelInfo))
		for k, v := range in.ModelInfo {
			out.ModelInfo[k] = ModelInfo{
				Name:    derefStr(v.Name),
				Version: derefStr(v.Version),
				Arch:    derefStr(v.Arch),
			}
		}
	}
	if len(in.Warnings) > 0 {
		out.Warnings = make([]Warning, len(in.Warnings))
		for i, w := range in.Warnings {
			out.Warnings[i] = Warning{
				Parameter: derefStr(w.Parameter),
				Type:      string(w.Type),
				Message:   derefStr(w.Message),
			}
		}
	}
	return out
}

func fromSpeechStarted(in *spectypes.SpeechStarted) *SpeechStartedEvent {
	return &SpeechStartedEvent{
		Channel:   int16SliceToInt(in.Channel),
		Timestamp: derefF32(in.Timestamp),
	}
}

func fromUtteranceEnd(in *spectypes.UtteranceEnd) *UtteranceEndEvent {
	return &UtteranceEndEvent{
		Channel:     int16SliceToInt(in.Channel),
		LastWordEnd: derefF32(in.LastWordEnd),
	}
}

func fromWsError(in *spectypes.WsError) *ErrorEvent {
	return &ErrorEvent{
		Variant:     string(in.Variant),
		Code:        derefStr(in.Code),
		Description: derefStr(in.Description),
		Message:     derefStr(in.Message),
	}
}

func fromServerSync(in *spectypes.ServerSync) *SyncEvent {
	return &SyncEvent{ID: derefInt64(in.Id)}
}
