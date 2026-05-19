// SAFETY-NET WIRE TESTS — DO NOT REMOVE WITHOUT EXPLICIT HUMAN DIRECTIVE.
//
// One test per public facade-options field that flows through
// optionsToTranscribeInput. These are the human-authored counterpart
// to wire_test_generated.go (which spec-codegen-go regenerates from
// the spec on every codegen run). A field that flows in both files
// gets a manual check here too; either file failing is a signal that
// the facade has drifted from api/.
//
// requireWired / requireDropped / isZeroForWire live in wire_helpers.go
// so they're visible to both this file and the generated file.

package restv1

import (
	"testing"
)

func TestWires_Callback(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Callback: "https://example.invalid/cb"})
	requireWired(t, in, "Callback")
}

func TestWires_CallbackMethod(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{CallbackMethod: "POST"})
	requireWired(t, in, "CallbackMethod")
}

func TestWires_DetectEntities(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{DetectEntities: true})
	requireWired(t, in, "DetectEntities")
}

func TestWires_DetectLanguage(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{DetectLanguage: []string{"en", "es"}})
	requireWired(t, in, "DetectLanguage")
}

func TestWires_Diarize(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Diarize: true})
	requireWired(t, in, "Diarize")
}

func TestWires_DiarizeModel(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{DiarizeModel: "nova-3"})
	requireWired(t, in, "DiarizeModel")
}

func TestWires_DiarizeVersion(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{DiarizeVersion: "2021-04-08"})
	requireWired(t, in, "DiarizeVersion")
}

func TestWires_Dictation(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Dictation: true})
	requireWired(t, in, "Dictation")
}

func TestWires_Encoding(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Encoding: "linear16"})
	requireWired(t, in, "Encoding")
}

func TestWires_FillerWords(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{FillerWords: true})
	requireWired(t, in, "FillerWords")
}

func TestWires_Intents(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Intents: true})
	requireWired(t, in, "Intents")
}

func TestWires_Keyterm(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Keyterm: []string{"deepgram"}})
	requireWired(t, in, "Keyterm")
}

func TestWires_Keywords(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Keywords: []string{"deepgram"}})
	requireWired(t, in, "Keywords")
}

func TestWires_Language(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Language: "en-US"})
	requireWired(t, in, "Language")
}

func TestWires_LogData(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{LogData: true})
	requireWired(t, in, "LogData")
}

func TestWires_Measurements(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Measurements: true})
	requireWired(t, in, "Measurements")
}

func TestWires_MipOptOut(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{MipOptOut: true})
	requireWired(t, in, "MipOptOut")
}

func TestWires_Model(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Model: "nova-3"})
	requireWired(t, in, "Model")
}

func TestWires_Multichannel(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Multichannel: true})
	requireWired(t, in, "Multichannel")
}

func TestWires_Numerals(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Numerals: true})
	requireWired(t, in, "Numerals")
}

func TestWires_Paragraphs(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Paragraphs: true})
	requireWired(t, in, "Paragraphs")
}

func TestWires_ProfanityFilter(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{ProfanityFilter: true})
	requireWired(t, in, "ProfanityFilter")
}

func TestWires_Punctuate(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Punctuate: true})
	requireWired(t, in, "Punctuate")
}

func TestWires_Redact(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Redact: []string{"pii"}})
	requireWired(t, in, "Redact")
}

func TestWires_Replace(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Replace: []string{"foo:bar"}})
	requireWired(t, in, "Replace")
}

func TestWires_Search(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Search: []string{"hello"}})
	requireWired(t, in, "Search")
}

func TestWires_Sentiment(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Sentiment: true})
	requireWired(t, in, "Sentiment")
}

func TestWires_SmartFormat(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{SmartFormat: true})
	requireWired(t, in, "SmartFormat")
}

func TestWires_Summarize(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Summarize: "v2"})
	requireWired(t, in, "Summarize")
}

func TestWires_Tag(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Tag: []string{"my-tag"}})
	requireWired(t, in, "Tag")
}

func TestWires_Topics(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Topics: true})
	requireWired(t, in, "Topics")
}

func TestWires_UttSplit(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{UttSplit: 0.8})
	requireWired(t, in, "UttSplit")
}

func TestWires_Utterances(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Utterances: true})
	requireWired(t, in, "Utterances")
}

func TestWires_VadTurnoff(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{VadTurnoff: 500})
	requireWired(t, in, "VadTurnoff")
}

func TestWires_Version(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Version: "latest"})
	requireWired(t, in, "Version")
}
