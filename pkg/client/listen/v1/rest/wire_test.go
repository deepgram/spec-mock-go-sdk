// SAFETY-NET WIRE TESTS — DO NOT REMOVE WITHOUT EXPLICIT HUMAN DIRECTIVE.
//
// These tests cover every field of PreRecordedTranscriptionOptions, so the
// whole facade is held to a wire-coverage invariant as the SDK gen builds
// out. There are two test families:
//
//   TestWires_<Field>   — facade field that SHOULD flow through to the
//                         generated wire type spectypes.TranscribeInput.
//                         Fails if the spec doesn't model the field yet,
//                         OR if the converter doesn't wire it. Drives the
//                         "this still needs work" signal during gen.
//
//   TestDropped_<Field> — facade field that is INTENTIONALLY not wired
//                         (deprecated server-side, never modeled in spec,
//                         or stem-internal). Documents the reason and
//                         asserts the converter leaves the corresponding
//                         spec field zero IF it exists.
//
// Each test is its own function so partial-deletion accidents stand out
// in code review. The shared helpers use reflection on the GENERATED type
// so a test compiles even when the spec doesn't model that field yet —
// failure surfaces as a test FAIL, not a compile error.
//
// Removal policy:
//   - LLM regens MUST NOT delete or edit any test in this file as a side
//     effect of an api/ change. If a facade field is genuinely removed
//     (compile-fenced by the literal on the test), removing the matching
//     test still requires a human directive in the triggering PR body.
//   - When a generated input field is renamed or retyped, update the
//     conventions in SKILL.md if needed and update assertions here, but
//     never remove a test.
//   - Adding a new test is REQUIRED whenever a new facade-options field
//     is added (whether it ends up wired or dropped).
//
// See `.agents/skills/sdk-facade-conventions/SKILL.md` for the conventions
// these tests enforce.

package restv1

import (
	"reflect"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces/v1"
)

// requireWired asserts that fieldName exists on spectypes.TranscribeInput
// AND that the converter set it to a non-zero value for the given opts.
// Uses reflection so the test compiles even when the spec hasn't modeled
// fieldName yet (it then fails at runtime with a clear message).
func requireWired(t *testing.T, in *spectypes.TranscribeInput, fieldName string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		t.Fatalf("spectypes.TranscribeInput has no field %q yet — spec needs to model it before the converter can wire it through. Until then this test SHOULD fail to keep the gap visible.", fieldName)
	}
	if isZeroForWire(v) {
		t.Fatalf("spectypes.TranscribeInput.%s exists but optionsToTranscribeInput didn't wire it. Add the corresponding `if o.%s ...` block (see sdk-facade-conventions/SKILL.md).", fieldName, fieldName)
	}
}

// requireDropped asserts the converter does NOT wire fieldName, which is
// the documented intent for facade-only fields with no spec counterpart
// (deprecated, stem-internal, etc.). If the spec later adds the field,
// the test still passes as long as the converter leaves it zero; promote
// to TestWires_<Field> at that point.
func requireDropped(t *testing.T, in *spectypes.TranscribeInput, fieldName, reason string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		return
	}
	if !isZeroForWire(v) {
		t.Fatalf("spectypes.TranscribeInput.%s is documented as permanently dropped (%s) but the converter wired it anyway. Either drop the wiring or update the rationale.", fieldName, reason)
	}
}

func isZeroForWire(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.String:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}

// =============================================================================
// TestWires_* — facade fields that should flow through to the wire
// =============================================================================

func TestWires_Alternatives(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Alternatives: 2})
	requireWired(t, in, "Alternatives")
}

func TestWires_Callback(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Callback: "https://example.invalid/cb"})
	requireWired(t, in, "Callback")
}

func TestWires_CallbackMethod(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{CallbackMethod: "POST"})
	requireWired(t, in, "CallbackMethod")
}

func TestWires_Channels(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Channels: 2})
	requireWired(t, in, "Channels")
}

func TestWires_DetectEntities(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{DetectEntities: true})
	requireWired(t, in, "DetectEntities")
}

func TestWires_DetectLanguage(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{DetectLanguage: true})
	requireWired(t, in, "DetectLanguage")
}

func TestWires_Diarize(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Diarize: true})
	requireWired(t, in, "Diarize")
}

func TestWires_DiarizeVersion(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{DiarizeVersion: "2025-01-01"})
	requireWired(t, in, "DiarizeVersion")
}

func TestWires_Dictation(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Dictation: true})
	requireWired(t, in, "Dictation")
}

func TestWires_Encoding(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Encoding: "linear16"})
	requireWired(t, in, "Encoding")
}

func TestWires_FillerWords(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{FillerWords: true})
	requireWired(t, in, "FillerWords")
}

func TestWires_Intents(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Intents: true})
	requireWired(t, in, "Intents")
}

func TestWires_Keywords(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Keywords: []string{"hello"}})
	requireWired(t, in, "Keywords")
}

func TestWires_Keyterm(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Keyterm: []string{"hello"}})
	requireWired(t, in, "Keyterm")
}

func TestWires_Language(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Language: "en-US"})
	requireWired(t, in, "Language")
}

func TestWires_Measurements(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Measurements: true})
	requireWired(t, in, "Measurements")
}

func TestWires_Model(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3"})
	requireWired(t, in, "Model")
}

func TestWires_Multichannel(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Multichannel: true})
	requireWired(t, in, "Multichannel")
}

func TestWires_Numerals(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Numerals: true})
	requireWired(t, in, "Numerals")
}

func TestWires_Paragraphs(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Paragraphs: true})
	requireWired(t, in, "Paragraphs")
}

func TestWires_ProfanityFilter(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{ProfanityFilter: true})
	requireWired(t, in, "ProfanityFilter")
}

func TestWires_Punctuate(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Punctuate: true})
	requireWired(t, in, "Punctuate")
}

func TestWires_Redact(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Redact: []string{"pci"}})
	requireWired(t, in, "Redact")
}

func TestWires_SampleRate(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{SampleRate: 16000})
	requireWired(t, in, "SampleRate")
}

func TestWires_Search(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Search: []string{"foo"}})
	requireWired(t, in, "Search")
}

func TestWires_Sentiment(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Sentiment: true})
	requireWired(t, in, "Sentiment")
}

func TestWires_SmartFormat(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{SmartFormat: true})
	requireWired(t, in, "SmartFormat")
}

func TestWires_Summarize(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Summarize: "v2"})
	requireWired(t, in, "Summarize")
}

func TestWires_Tag(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Tag: []string{"prod"}})
	requireWired(t, in, "Tag")
}

func TestWires_Topics(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Topics: true})
	requireWired(t, in, "Topics")
}

func TestWires_UttSplit(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{UttSplit: 1.5})
	requireWired(t, in, "UttSplit")
}

func TestWires_Utterances(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Utterances: true})
	requireWired(t, in, "Utterances")
}

func TestWires_Version(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Version: "2025-01-01"})
	requireWired(t, in, "Version")
}

// =============================================================================
// TestDropped_* — facade fields intentionally not wired
// =============================================================================

func TestDropped_CustomIntent(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{CustomIntent: []string{"x"}})
	requireDropped(t, in, "CustomIntent", "Stem ingests this server-side via /v1/listen but the spec doesn't model the field; no plans to expose through the generated wire type.")
}

func TestDropped_CustomIntentMode(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{CustomIntentMode: "extended"})
	requireDropped(t, in, "CustomIntentMode", "Companion to CustomIntent; same status — server-side only.")
}

func TestDropped_CustomTopic(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{CustomTopic: []string{"x"}})
	requireDropped(t, in, "CustomTopic", "Stem ingests this server-side via /v1/listen but the spec doesn't model the field; no plans to expose through the generated wire type.")
}

func TestDropped_CustomTopicMode(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{CustomTopicMode: "extended"})
	requireDropped(t, in, "CustomTopicMode", "Companion to CustomTopic; same status — server-side only.")
}

func TestDropped_DetectTopics(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{DetectTopics: true})
	requireDropped(t, in, "DetectTopics", "Fully deprecated by stem 12/25 in favour of `topics`. The spec correctly omits it; the facade keeps it for source-compat with deepgram-go-sdk callers.")
}

func TestDropped_Extra(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Extra: []string{"x=y"}})
	requireDropped(t, in, "Extra", "Stem-side metadata pass-through; the spec models the response-side `extra` ExtraMap on ResponseMetadata, not the request side.")
}

func TestDropped_Replace(t *testing.T) {
	in := optionsToTranscribeInput(&interfaces.PreRecordedTranscriptionOptions{Replace: []string{"foo:bar"}})
	requireDropped(t, in, "Replace", "Stem skips this server-side; no plans to add to spec.")
}
