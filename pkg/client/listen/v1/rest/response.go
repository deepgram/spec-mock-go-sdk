// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Idiomatic value-typed response shapes for the listen REST surface.
// Customers receive these from FromURL / FromFile / FromStream. The
// facade absorbs the structural differences between successive smithy-go
// regenerations of api/types - field renames, nullability changes, and
// pointer vs value typing - so call sites stay stable across spec
// evolution. Conversions from the generated spectypes.TranscribeOutput
// to these shapes live in convert.go.

package restv1

type SummaryInfo struct {
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
	ModelUUID    string `json:"model_uuid,omitempty"`
}

type ModelInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Arch    string `json:"arch,omitempty"`
}

type IntentsInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type SentimentInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type TopicsInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type Metadata struct {
	TransactionKey string               `json:"transaction_key,omitempty"`
	RequestID      string               `json:"request_id,omitempty"`
	Sha256         string               `json:"sha256,omitempty"`
	Created        string               `json:"created,omitempty"`
	Duration       float64              `json:"duration,omitempty"`
	Channels       int                  `json:"channels,omitempty"`
	Models         []string             `json:"models,omitempty"`
	ModelInfo      map[string]ModelInfo `json:"model_info,omitempty"`
	Warnings       *[]Warning           `json:"warnings,omitempty"`
	SummaryInfo    *SummaryInfo         `json:"summary_info,omitempty"`
	IntentsInfo    *IntentsInfo         `json:"intents_info,omitempty"`
	SentimentInfo  *SentimentInfo       `json:"sentiment_info,omitempty"`
	TopicsInfo     *TopicsInfo          `json:"topics_info,omitempty"`
	Extra          map[string]string    `json:"extra,omitempty"`
}

type Warning struct {
	Parameter string `json:"parameter,omitempty"`
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
}

type Hit struct {
	Confidence float64 `json:"confidence,omitempty"`
	Start      float64 `json:"start,omitempty"`
	End        float64 `json:"end,omitempty"`
	Snippet    string  `json:"snippet,omitempty"`
}

type Search struct {
	Query string `json:"query,omitempty"`
	Hits  []Hit  `json:"hits,omitempty"`
}

type Word struct {
	Word              string   `json:"word,omitempty"`
	Start             float64  `json:"start,omitempty"`
	End               float64  `json:"end,omitempty"`
	Confidence        float64  `json:"confidence,omitempty"`
	Speaker           *int     `json:"speaker,omitempty"`
	SpeakerConfidence *float64 `json:"speaker_confidence,omitempty"`
	PunctuatedWord    string   `json:"punctuated_word,omitempty"`
	Sentiment         *string  `json:"sentiment,omitempty"`
	SentimentScore    *float64 `json:"sentiment_score,omitempty"`
	Language          string   `json:"language,omitempty"`
}

type Translation struct {
	Language    string `json:"language,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type Alternative struct {
	Transcript  string       `json:"transcript,omitempty"`
	Confidence  float64      `json:"confidence,omitempty"`
	Words       []Word       `json:"words,omitempty"`
	Paragraphs  *Paragraphs  `json:"paragraphs,omitempty"`
	Entities    *[]Entity    `json:"entities,omitempty"`
	Summaries   *[]SummaryV1 `json:"summaries,omitempty"`
	Translation *Translation `json:"translation,omitempty"`
	Languages   []string     `json:"languages,omitempty"`
}

type Paragraphs struct {
	Transcript string      `json:"transcript,omitempty"`
	Paragraphs []Paragraph `json:"paragraphs,omitempty"`
}

type Paragraph struct {
	Sentences      []Sentence `json:"sentences,omitempty"`
	NumWords       int        `json:"num_words,omitempty"`
	Start          float64    `json:"start,omitempty"`
	End            float64    `json:"end,omitempty"`
	Speaker        *int       `json:"speaker,omitempty"`
	Sentiment      *string    `json:"sentiment,omitempty"`
	SentimentScore *float64   `json:"sentiment_score,omitempty"`
}

type Sentence struct {
	Text           string   `json:"text,omitempty"`
	Start          float64  `json:"start,omitempty"`
	End            float64  `json:"end,omitempty"`
	Sentiment      *string  `json:"sentiment,omitempty"`
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
}

type Entity struct {
	Label      string  `json:"label,omitempty"`
	Value      string  `json:"value,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	StartWord  float64 `json:"start_word,omitempty"`
	EndWord    float64 `json:"end_word,omitempty"`
}

type Channel struct {
	Search             *[]Search     `json:"search,omitempty"`
	Alternatives       []Alternative `json:"alternatives,omitempty"`
	DetectedLanguage   string        `json:"detected_language,omitempty"`
	LanguageConfidence float64       `json:"language_confidence,omitempty"`
}

type Utterance struct {
	Start          float64  `json:"start,omitempty"`
	End            float64  `json:"end,omitempty"`
	Confidence     float64  `json:"confidence,omitempty"`
	Channel        int      `json:"channel,omitempty"`
	Transcript     string   `json:"transcript,omitempty"`
	Words          []Word   `json:"words,omitempty"`
	Speaker        *int     `json:"speaker,omitempty"`
	Sentiment      *string  `json:"sentiment,omitempty"`
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
	ID             string   `json:"id,omitempty"`
}

type SummaryV1 struct {
	Summary   string `json:"summary,omitempty"`
	StartWord int    `json:"start_word,omitempty"`
	EndWord   int    `json:"end_word,omitempty"`
}

type Result struct {
	Channels   []Channel   `json:"channels,omitempty"`
	Utterances []Utterance `json:"utterances,omitempty"`
}

// PreRecordedResponse is the customer-facing response from FromURL /
// FromFile / FromStream. Stable contract regardless of how api/types
// evolves under regeneration.
type PreRecordedResponse struct {
	RequestID string    `json:"request_id,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
	Results   *Result   `json:"results,omitempty"`
}
