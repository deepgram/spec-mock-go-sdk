package restv1

// PreRecordedTranscriptionOptions are the customer-facing options for
// POST /v1/listen. Each exported field corresponds to a single
// @httpQuery member on the generated spectypes.TranscribeInput wire
// shape; the client exposes idiomatic Go value types so callers do
// not nil-check primitives.
//
// Zero values mean "not set" and leave the corresponding wire field
// nil so the parameter does not appear in the request query string.
//
// The set of fields here is the public surface of /v1/listen. Any
// @internal-tagged parameter on the wire input shape is intentionally
// absent.
type PreRecordedTranscriptionOptions struct {
	Callback        string   `json:"callback,omitempty"         schema:"callback,omitempty"`
	CallbackMethod  string   `json:"callback_method,omitempty"  schema:"callback_method,omitempty"`
	DetectEntities  bool     `json:"detect_entities,omitempty"  schema:"detect_entities,omitempty"`
	DetectLanguage  []string `json:"detect_language,omitempty"  schema:"detect_language,omitempty"`

	// Deprecated: Legacy flag. Prefer DiarizeModel for explicit model selection.
	// Diarize=true continues to work for backward compatibility but is mutually
	// exclusive with DiarizeModel.
	Diarize      bool   `json:"diarize,omitempty"       schema:"diarize,omitempty"`
	DiarizeModel string `json:"diarize_model,omitempty" schema:"diarize_model,omitempty"`

	// Deprecated: Legacy diarization-model selector. Prefer DiarizeModel.
	// Mutually exclusive with DiarizeModel.
	DiarizeVersion string `json:"diarize_version,omitempty" schema:"diarize_version,omitempty"`

	Dictation   bool     `json:"dictation,omitempty"    schema:"dictation,omitempty"`
	Encoding    string   `json:"encoding,omitempty"     schema:"encoding,omitempty"`
	FillerWords bool     `json:"filler_words,omitempty" schema:"filler_words,omitempty"`
	Intents     bool     `json:"intents,omitempty"      schema:"intents,omitempty"`
	Keyterm     []string `json:"keyterm,omitempty"      schema:"keyterm,omitempty"`
	Keywords    []string `json:"keywords,omitempty"     schema:"keywords,omitempty"`
	Language    string   `json:"language,omitempty"     schema:"language,omitempty"`

	// Deprecated: Prefer MipOptOut. LogData is recognized for backward
	// compatibility; sending both with conflicting values returns 400.
	LogData bool `json:"log_data,omitempty" schema:"log_data,omitempty"`

	Measurements    bool     `json:"measurements,omitempty"     schema:"measurements,omitempty"`
	MipOptOut       bool     `json:"mip_opt_out,omitempty"      schema:"mip_opt_out,omitempty"`
	Model           string   `json:"model,omitempty"            schema:"model,omitempty"`
	Multichannel    bool     `json:"multichannel,omitempty"     schema:"multichannel,omitempty"`
	Numerals        bool     `json:"numerals,omitempty"         schema:"numerals,omitempty"`
	Paragraphs      bool     `json:"paragraphs,omitempty"       schema:"paragraphs,omitempty"`
	ProfanityFilter bool     `json:"profanity_filter,omitempty" schema:"profanity_filter,omitempty"`
	Punctuate       bool     `json:"punctuate,omitempty"        schema:"punctuate,omitempty"`
	Redact          []string `json:"redact,omitempty"           schema:"redact,omitempty"`
	Replace         []string `json:"replace,omitempty"          schema:"replace,omitempty"`
	Search          []string `json:"search,omitempty"           schema:"search,omitempty"`
	Sentiment       bool     `json:"sentiment,omitempty"        schema:"sentiment,omitempty"`
	SmartFormat     bool     `json:"smart_format,omitempty"     schema:"smart_format,omitempty"`
	Summarize       string   `json:"summarize,omitempty"        schema:"summarize,omitempty"`
	Tag             []string `json:"tag,omitempty"              schema:"tag,omitempty"`
	Topics          bool     `json:"topics,omitempty"           schema:"topics,omitempty"`
	UttSplit        float64  `json:"utt_split,omitempty"        schema:"utt_split,omitempty"`
	Utterances      bool     `json:"utterances,omitempty"       schema:"utterances,omitempty"`

	// Deprecated: Use Endpointing instead. VadTurnoff is rejected when
	// Endpointing is also set as an integer; if only VadTurnoff is set it
	// is silently mapped onto Endpointing.
	VadTurnoff int `json:"vad_turnoff,omitempty" schema:"vad_turnoff,omitempty"`

	Version string `json:"version,omitempty" schema:"version,omitempty"`
}
