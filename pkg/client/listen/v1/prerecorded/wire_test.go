// SAFETY-NET WIRE TESTS — DO NOT REMOVE WITHOUT EXPLICIT HUMAN DIRECTIVE.

package prerecordedv1

import (
	"reflect"
	"testing"
)

func TestWires_Callback(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Callback: "https://example.invalid/callback"})
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
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Keywords: []string{"deepgram:2"}})
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
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Replace: []string{"hello:hi"}})
	requireWired(t, in, "Replace")
}

func TestWires_Search(t *testing.T) {
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Search: []string{"project"}})
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
	in := optionsToTranscribeInput(&PreRecordedTranscriptionOptions{Tag: []string{"demo"}})
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

func TestFacadeOnly_AdditionalQueryParams(t *testing.T) {
	requireFacadeOnly(t, &PreRecordedTranscriptionOptions{}, "AdditionalQueryParams")
}

func requireWired(t *testing.T, input any, field string) {
	t.Helper()
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		t.Fatalf("input must be a non-nil pointer, got %T", input)
	}
	f := v.Elem().FieldByName(field)
	if !f.IsValid() {
		t.Fatalf("input %T has no field %q", input, field)
	}
	if isZeroForWire(f) {
		t.Fatalf("input %T field %q was not wired", input, field)
	}
}

func requireDropped(t *testing.T, input any, field string) {
	t.Helper()
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		t.Fatalf("input must be a non-nil pointer, got %T", input)
	}
	f := v.Elem().FieldByName(field)
	if !f.IsValid() {
		return
	}
	if !isZeroForWire(f) {
		t.Fatalf("input %T field %q should not be wired", input, field)
	}
}

func requireFacadeOnly(t *testing.T, opts any, field string) {
	t.Helper()
	opt := reflect.ValueOf(opts)
	if opt.Kind() != reflect.Ptr || opt.IsNil() {
		t.Fatalf("opts must be a non-nil pointer, got %T", opts)
	}
	if !opt.Elem().FieldByName(field).IsValid() {
		t.Fatalf("options %T has no field %q", opts, field)
	}
	in := optionsToTranscribeInput(opts.(*PreRecordedTranscriptionOptions))
	if reflect.ValueOf(in).Elem().FieldByName(field).IsValid() {
		t.Fatalf("wire input unexpectedly has field %q", field)
	}
}

func isZeroForWire(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.String:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}
