// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1

import (
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

// optionsToTranscribeInput translates an idiomatic
// *PreRecordedTranscriptionOptions into a *spectypes.TranscribeInput
// suitable for httptransport.Invoke. Zero-valued facade fields leave the
// generated pointer as nil so the parameter is omitted from the wire
// query string.
//
// Every @httpQuery member of spectypes.TranscribeInput has a wiring
// block here. The set of @internal-tagged parameters that the spec
// excludes from the public surface is intentionally not represented on
// PreRecordedTranscriptionOptions and therefore has no wiring block.
func optionsToTranscribeInput(o *PreRecordedTranscriptionOptions) *spectypes.TranscribeInput {
	in := &spectypes.TranscribeInput{}
	if o == nil {
		return in
	}
	if o.Callback != "" {
		v := o.Callback
		in.Callback = &v
	}
	if o.CallbackMethod != "" {
		in.CallbackMethod = spectypes.CallbackMethod(o.CallbackMethod)
	}
	if o.DetectEntities {
		v := o.DetectEntities
		in.DetectEntities = &v
	}
	if len(o.DetectLanguage) > 0 {
		in.DetectLanguage = o.DetectLanguage
	}
	if o.Diarize {
		v := o.Diarize
		in.Diarize = &v
	}
	if o.DiarizeModel != "" {
		v := o.DiarizeModel
		in.DiarizeModel = &v
	}
	if o.DiarizeVersion != "" {
		v := o.DiarizeVersion
		in.DiarizeVersion = &v
	}
	if o.Dictation {
		v := o.Dictation
		in.Dictation = &v
	}
	if o.Encoding != "" {
		v := o.Encoding
		in.Encoding = &v
	}
	if o.FillerWords {
		v := o.FillerWords
		in.FillerWords = &v
	}
	if o.Intents {
		v := o.Intents
		in.Intents = &v
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
	if o.Measurements {
		v := o.Measurements
		in.Measurements = &v
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
	if o.Numerals {
		v := o.Numerals
		in.Numerals = &v
	}
	if o.Paragraphs {
		v := o.Paragraphs
		in.Paragraphs = &v
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
	if len(o.Replace) > 0 {
		in.Replace = o.Replace
	}
	if len(o.Search) > 0 {
		in.Search = o.Search
	}
	if o.Sentiment {
		v := o.Sentiment
		in.Sentiment = &v
	}
	if o.SmartFormat {
		v := o.SmartFormat
		in.SmartFormat = &v
	}
	if o.Summarize != "" {
		v := o.Summarize
		in.Summarize = &v
	}
	if len(o.Tag) > 0 {
		in.Tag = o.Tag
	}
	if o.Topics {
		v := o.Topics
		in.Topics = &v
	}
	if o.UttSplit != 0 {
		v := float32(o.UttSplit)
		in.UttSplit = &v
	}
	if o.Utterances {
		v := o.Utterances
		in.Utterances = &v
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

func derefInt64(p *int64) int {
	if p == nil {
		return 0
	}
	return int(*p)
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

func sentimentToPtrStr(s spectypes.Sentiment) *string {
	if s == "" {
		return nil
	}
	v := string(s)
	return &v
}

func convertTranscribeOutput(in *spectypes.TranscribeOutput) *PreRecordedResponse {
	if in == nil {
		return nil
	}
	return &PreRecordedResponse{
		RequestID: derefStr(in.RequestId),
		Metadata:  convertResponseMetadata(in.Metadata),
		Results:   convertResponseResults(in.Results),
	}
}

func convertResponseMetadata(in *spectypes.ResponseMetadata) *Metadata {
	if in == nil {
		return nil
	}
	out := &Metadata{
		TransactionKey: derefStr(in.TransactionKey),
		RequestID:      derefStr(in.RequestId),
		Sha256:         derefStr(in.Sha256),
		Created:        derefStr(in.Created),
		Duration:       derefF32(in.Duration),
		Channels:       derefInt32(in.Channels),
		Extra:          in.Extra,
	}
	if len(in.Models) > 0 {
		out.Models = in.Models
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
		ws := make([]Warning, len(in.Warnings))
		for i, w := range in.Warnings {
			ws[i] = Warning{
				Parameter: derefStr(w.Parameter),
				Type:      string(w.Type),
				Message:   derefStr(w.Message),
			}
		}
		out.Warnings = &ws
	}
	if in.SummaryInfo != nil {
		out.SummaryInfo = &SummaryInfo{
			InputTokens:  derefInt64(in.SummaryInfo.InputTokens),
			OutputTokens: derefInt64(in.SummaryInfo.OutputTokens),
			ModelUUID:    derefStr(in.SummaryInfo.ModelUuid),
		}
	}
	if in.IntentsInfo != nil {
		out.IntentsInfo = &IntentsInfo{
			InputTokens:  derefInt64(in.IntentsInfo.InputTokens),
			OutputTokens: derefInt64(in.IntentsInfo.OutputTokens),
			ModelUUID:    derefStr(in.IntentsInfo.ModelUuid),
		}
	}
	if in.SentimentInfo != nil {
		out.SentimentInfo = &SentimentInfo{
			InputTokens:  derefInt64(in.SentimentInfo.InputTokens),
			OutputTokens: derefInt64(in.SentimentInfo.OutputTokens),
			ModelUUID:    derefStr(in.SentimentInfo.ModelUuid),
		}
	}
	if in.TopicsInfo != nil {
		out.TopicsInfo = &TopicsInfo{
			InputTokens:  derefInt64(in.TopicsInfo.InputTokens),
			OutputTokens: derefInt64(in.TopicsInfo.OutputTokens),
			ModelUUID:    derefStr(in.TopicsInfo.ModelUuid),
		}
	}
	return out
}

func convertResponseResults(in *spectypes.ResponseResults) *Result {
	if in == nil {
		return nil
	}
	out := &Result{}
	if len(in.Channels) > 0 {
		out.Channels = make([]Channel, len(in.Channels))
		for i, c := range in.Channels {
			out.Channels[i] = convertChannel(&c)
		}
	}
	if len(in.Utterances) > 0 {
		out.Utterances = make([]Utterance, len(in.Utterances))
		for i, u := range in.Utterances {
			out.Utterances[i] = convertUtterance(&u)
		}
	}
	return out
}

func convertChannel(in *spectypes.Channel) Channel {
	out := Channel{
		DetectedLanguage:   derefStr(in.DetectedLanguage),
		LanguageConfidence: derefF32(in.LanguageConfidence),
	}
	if len(in.Search) > 0 {
		ss := make([]Search, len(in.Search))
		for i, s := range in.Search {
			ss[i] = convertSearch(&s)
		}
		out.Search = &ss
	}
	if len(in.Alternatives) > 0 {
		out.Alternatives = make([]Alternative, len(in.Alternatives))
		for i, a := range in.Alternatives {
			out.Alternatives[i] = convertAlternative(&a)
		}
	}
	return out
}

func convertSearch(in *spectypes.SearchResult) Search {
	out := Search{
		Query: derefStr(in.Query),
	}
	if len(in.Hits) > 0 {
		out.Hits = make([]Hit, len(in.Hits))
		for i, h := range in.Hits {
			out.Hits[i] = Hit{
				Confidence: derefF32(h.Confidence),
				Start:      derefF32(h.Start),
				End:        derefF32(h.End),
				Snippet:    derefStr(h.Snippet),
			}
		}
	}
	return out
}

func convertAlternative(in *spectypes.Alternative) Alternative {
	out := Alternative{
		Transcript: derefStr(in.Transcript),
		Confidence: derefF32(in.Confidence),
		Languages:  in.Languages,
	}
	if len(in.Words) > 0 {
		out.Words = make([]Word, len(in.Words))
		for i, w := range in.Words {
			out.Words[i] = convertWord(&w)
		}
	}
	if in.Paragraphs != nil {
		out.Paragraphs = convertParagraphs(in.Paragraphs)
	}
	if len(in.Entities) > 0 {
		es := make([]Entity, len(in.Entities))
		for i, e := range in.Entities {
			es[i] = Entity{
				Label:      derefStr(e.Label),
				Value:      derefStr(e.Value),
				Confidence: derefF32(e.Confidence),
				StartWord:  float64(derefInt32(e.StartWord)),
				EndWord:    float64(derefInt32(e.EndWord)),
			}
		}
		out.Entities = &es
	}
	if len(in.Summaries) > 0 {
		ss := make([]SummaryV1, len(in.Summaries))
		for i, s := range in.Summaries {
			ss[i] = SummaryV1{
				Summary:   derefStr(s.Summary),
				StartWord: derefInt32(s.StartWord),
				EndWord:   derefInt32(s.EndWord),
			}
		}
		out.Summaries = &ss
	}
	return out
}

func convertWord(in *spectypes.Word) Word {
	return Word{
		Word:              derefStr(in.Word),
		Start:             derefF32(in.Start),
		End:               derefF32(in.End),
		Confidence:        derefF32(in.Confidence),
		Speaker:           int64ToPtrInt(in.Speaker),
		SpeakerConfidence: f32ToPtrF64(in.SpeakerConfidence),
		PunctuatedWord:    derefStr(in.PunctuatedWord),
		Sentiment:         sentimentToPtrStr(in.Sentiment),
		SentimentScore:    f32ToPtrF64(in.SentimentScore),
		Language:          derefStr(in.Language),
	}
}

func convertParagraphs(in *spectypes.Paragraphs) *Paragraphs {
	out := &Paragraphs{
		Transcript: derefStr(in.Transcript),
	}
	if len(in.Paragraphs) > 0 {
		out.Paragraphs = make([]Paragraph, len(in.Paragraphs))
		for i, p := range in.Paragraphs {
			out.Paragraphs[i] = Paragraph{
				NumWords:       derefInt32(p.NumWords),
				Start:          derefF32(p.Start),
				End:            derefF32(p.End),
				Speaker:        int64ToPtrInt(p.Speaker),
				Sentiment:      sentimentToPtrStr(p.Sentiment),
				SentimentScore: f32ToPtrF64(p.SentimentScore),
			}
			if len(p.Sentences) > 0 {
				out.Paragraphs[i].Sentences = make([]Sentence, len(p.Sentences))
				for j, s := range p.Sentences {
					out.Paragraphs[i].Sentences[j] = Sentence{
						Text:           derefStr(s.Text),
						Start:          derefF32(s.Start),
						End:            derefF32(s.End),
						Sentiment:      sentimentToPtrStr(s.Sentiment),
						SentimentScore: f32ToPtrF64(s.SentimentScore),
					}
				}
			}
		}
	}
	return out
}

func convertUtterance(in *spectypes.Utterance) Utterance {
	out := Utterance{
		Start:          derefF32(in.Start),
		End:            derefF32(in.End),
		Confidence:     derefF32(in.Confidence),
		Channel:        derefInt32(in.Channel),
		Transcript:     derefStr(in.Transcript),
		Speaker:        int64ToPtrInt(in.Speaker),
		Sentiment:      sentimentToPtrStr(in.Sentiment),
		SentimentScore: f32ToPtrF64(in.SentimentScore),
		ID:             derefStr(in.Id),
	}
	if len(in.Words) > 0 {
		out.Words = make([]Word, len(in.Words))
		for i, w := range in.Words {
			out.Words[i] = convertWord(&w)
		}
	}
	return out
}
