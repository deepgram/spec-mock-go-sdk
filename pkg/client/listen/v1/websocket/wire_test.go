// SAFETY-NET WIRE TESTS — DO NOT REMOVE WITHOUT EXPLICIT HUMAN DIRECTIVE.
//
// One test per public facade-options field that flows through
// optionsToStreamInput. spec-codegen-go does not currently generate
// streaming wire-test stubs (no @facadeOptionsType on the streaming
// operation), so this is the entire wire-test surface for live
// transcription. requireWired / requireDropped / isZeroForWire live
// in wire_helpers.go.

package wsv1

import (
	"testing"
)

func TestWires_Channels(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Channels: 2})
	requireWired(t, in, "Channels")
}

func TestWires_Diarize(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Diarize: true})
	requireWired(t, in, "Diarize")
}

func TestWires_DiarizeVersion(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{DiarizeVersion: "2021-04-08"})
	requireWired(t, in, "DiarizeVersion")
}

func TestWires_Encoding(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Encoding: "linear16"})
	requireWired(t, in, "Encoding")
}

func TestWires_Endpointing(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Endpointing: 500})
	requireWired(t, in, "Endpointing")
}

func TestWires_InterimResults(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{InterimResults: true})
	requireWired(t, in, "InterimResults")
}

func TestWires_Keyterm(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Keyterm: []string{"deepgram"}})
	requireWired(t, in, "Keyterm")
}

func TestWires_Keywords(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Keywords: []string{"deepgram"}})
	requireWired(t, in, "Keywords")
}

func TestWires_Language(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Language: "en-US"})
	requireWired(t, in, "Language")
}

func TestWires_LogData(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{LogData: true})
	requireWired(t, in, "LogData")
}

func TestWires_MipOptOut(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{MipOptOut: true})
	requireWired(t, in, "MipOptOut")
}

func TestWires_Model(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Model: "nova-3"})
	requireWired(t, in, "Model")
}

func TestWires_Multichannel(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Multichannel: true})
	requireWired(t, in, "Multichannel")
}

func TestWires_ProfanityFilter(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{ProfanityFilter: true})
	requireWired(t, in, "ProfanityFilter")
}

func TestWires_Punctuate(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Punctuate: true})
	requireWired(t, in, "Punctuate")
}

func TestWires_Redact(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Redact: []string{"pii"}})
	requireWired(t, in, "Redact")
}

func TestWires_SampleRate(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{SampleRate: 16000})
	requireWired(t, in, "SampleRate")
}

func TestWires_Search(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Search: []string{"hello"}})
	requireWired(t, in, "Search")
}

func TestWires_SmartFormat(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{SmartFormat: true})
	requireWired(t, in, "SmartFormat")
}

func TestWires_Tag(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Tag: []string{"my-tag"}})
	requireWired(t, in, "Tag")
}

func TestWires_UtteranceEndMs(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{UtteranceEndMs: 1000})
	requireWired(t, in, "UtteranceEndMs")
}

func TestWires_VadEvents(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{VadEvents: true})
	requireWired(t, in, "VadEvents")
}

func TestWires_VadTurnoff(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{VadTurnoff: 500})
	requireWired(t, in, "VadTurnoff")
}

func TestWires_Version(t *testing.T) {
	in := optionsToStreamInput(&LiveTranscriptionOptions{Version: "latest"})
	requireWired(t, in, "Version")
}
