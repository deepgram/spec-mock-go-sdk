// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Converters: spectypes.TranscribeOutput (generated, pointer-typed, in
// api/types) -> *PreRecordedResponse (idiomatic, value-typed, in this
// package). One pure-Go function per type pair. Nil pointer in =>
// zero-value field out, so customer code can read response fields
// without ever nil-checking.

package restv1

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

// convertTranscribeOutput converts the generated TranscribeOutput into
// the customer-facing PreRecordedResponse.
//
// NOTE: The generated TranscribeOutput.Metadata field was removed in
// this api/types regeneration. The customer-facing
// PreRecordedResponse.Metadata field is preserved (always nil) so the
// Go signature stays stable. See BREAKING_CHANGES.md for impact on
// customer code that read resp.Metadata.
func convertTranscribeOutput(in *spectypes.TranscribeOutput) *PreRecordedResponse {
	if in == nil {
		return nil
	}
	return &PreRecordedResponse{
		RequestID: derefStr(in.RequestId),
		Metadata:  nil, // FIELD_REMOVED on spectypes.TranscribeOutput; absorbed-with-ceremony
		Results:   convertResponseResults(in.Results),
	}
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
