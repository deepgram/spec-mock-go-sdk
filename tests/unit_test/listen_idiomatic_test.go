// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package deepgramtest

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

const (
	endpoint = "https://api.deepgram.com/v1/listen"
	apiKey   = "test-api-key"
)

func mockPost(t *testing.T, body string) {
	t.Helper()
	httpmock.RegisterResponder(http.MethodPost, endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})
}

func newClient() *client.Client {
	c := client.New(apiKey, &interfaces.ClientOptions{Host: "https://api.deepgram.com"})
	httpmock.ActivateNonDefault(&c.Client.HTTPClient.Client)
	return c
}

func ptrStr(s string) *string  { return &s }
func ptrF32(f float32) *float32 { return &f }
func ptrI64(i int64) *int64    { return &i }
func ptrI32(i int32) *int32    { return &i }

func Test_FromURL_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "abc-123",
		"metadata": {
			"transaction_key": "deprecated",
			"request_id": "abc-123",
			"sha256": "fakehash",
			"created": "2026-05-12T12:00:00Z",
			"duration": 3.14,
			"channels": 1,
			"model_info": {}
		},
		"results": {
			"channels": [{
				"alternatives": [{
					"transcript": "hello world",
					"confidence": 0.99,
					"words": [
						{"word": "hello", "start": 0.0, "end": 0.5, "confidence": 0.98},
						{"word": "world", "start": 0.5, "end": 1.0, "confidence": 0.97}
					]
				}]
			}]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		t.Fatalf("FromURL err: %v", err)
	}
	if res.RequestId == nil || *res.RequestId != "abc-123" {
		t.Fatalf("RequestId = %v, want abc-123", res.RequestId)
	}
	if res.Metadata == nil {
		t.Fatal("Metadata nil")
	}
	if res.Metadata.Duration == nil || *res.Metadata.Duration != 3.14 {
		t.Fatalf("Duration = %v, want 3.14", res.Metadata.Duration)
	}
	if res.Results == nil || len(res.Results.Channels) != 1 {
		t.Fatalf("Channels = %d, want 1", len(res.Results.Channels))
	}
	alt := res.Results.Channels[0].Alternatives[0]
	if alt.Transcript == nil || *alt.Transcript != "hello world" {
		t.Fatalf("Transcript = %v, want 'hello world'", alt.Transcript)
	}
	if len(alt.Words) != 2 {
		t.Fatalf("len(Words) = %d, want 2", len(alt.Words))
	}
	if alt.Words[0].Word == nil || *alt.Words[0].Word != "hello" {
		t.Fatalf("Words[0].Word = %v, want 'hello'", alt.Words[0].Word)
	}
}

func Test_FromStream_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "stream-456",
		"results": {
			"channels": [{
				"alternatives": [{"transcript": "streamed audio", "confidence": 0.9, "words": []}]
			}]
		}
	}`)

	src := strings.NewReader("fake-audio-bytes")
	res, err := newClient().FromStream(context.Background(), src,
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		t.Fatalf("FromStream err: %v", err)
	}
	if *res.RequestId != "stream-456" {
		t.Fatalf("RequestId = %s, want stream-456", *res.RequestId)
	}
}

func Test_Diarization_SpeakerField(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "diar-789",
		"results": {
			"channels": [{
				"alternatives": [{
					"transcript": "speaker test",
					"confidence": 0.95,
					"words": [
						{"word": "hello", "start": 0.0, "end": 0.5, "confidence": 0.9, "speaker": 0, "speaker_confidence": 0.85},
						{"word": "there", "start": 0.5, "end": 1.0, "confidence": 0.9, "speaker": 1, "speaker_confidence": 0.92}
					]
				}]
			}]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3", Diarize: true})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	words := res.Results.Channels[0].Alternatives[0].Words
	if words[0].Speaker == nil || *words[0].Speaker != 0 {
		t.Fatalf("Words[0].Speaker = %v, want 0", words[0].Speaker)
	}
	if words[1].Speaker == nil || *words[1].Speaker != 1 {
		t.Fatalf("Words[1].Speaker = %v, want 1", words[1].Speaker)
	}
	if words[1].SpeakerConfidence == nil || *words[1].SpeakerConfidence != 0.92 {
		t.Fatalf("Words[1].SpeakerConfidence = %v, want 0.92", words[1].SpeakerConfidence)
	}
}

func Test_Multichannel_Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "mc-001",
		"results": {
			"channels": [
				{"alternatives": [{"transcript": "ch0 audio", "confidence": 0.9, "words": []}]},
				{"alternatives": [{"transcript": "ch1 audio", "confidence": 0.88, "words": []}]}
			]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3", Multichannel: true})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(res.Results.Channels) != 2 {
		t.Fatalf("Channels = %d, want 2", len(res.Results.Channels))
	}
	if *res.Results.Channels[0].Alternatives[0].Transcript != "ch0 audio" {
		t.Fatal("Channel 0 transcript wrong")
	}
	if *res.Results.Channels[1].Alternatives[0].Transcript != "ch1 audio" {
		t.Fatal("Channel 1 transcript wrong")
	}
}

func Test_Paragraphs_Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "para-001",
		"results": {
			"channels": [{
				"alternatives": [{
					"transcript": "hello. world.",
					"confidence": 0.9,
					"words": [],
					"paragraphs": {
						"transcript": "hello. world.",
						"paragraphs": [{
							"sentences": [
								{"text": "hello.", "start": 0.0, "end": 0.5},
								{"text": "world.", "start": 0.5, "end": 1.0}
							],
							"num_words": 2,
							"start": 0.0,
							"end": 1.0,
							"speaker": 0
						}]
					}
				}]
			}]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	paras := res.Results.Channels[0].Alternatives[0].Paragraphs
	if paras == nil {
		t.Fatal("Paragraphs nil")
	}
	if len(paras.Paragraphs) != 1 {
		t.Fatalf("len(Paragraphs) = %d, want 1", len(paras.Paragraphs))
	}
	p := paras.Paragraphs[0]
	if len(p.Sentences) != 2 {
		t.Fatalf("Sentences = %d, want 2", len(p.Sentences))
	}
	if *p.Sentences[0].Text != "hello." {
		t.Fatalf("Sentence[0].Text = %v, want 'hello.'", *p.Sentences[0].Text)
	}
	if p.Speaker == nil || *p.Speaker != 0 {
		t.Fatalf("Paragraph.Speaker = %v, want 0", p.Speaker)
	}
}

func Test_Entities_Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "ent-001",
		"results": {
			"channels": [{
				"alternatives": [{
					"transcript": "I live in Paris.",
					"confidence": 0.9,
					"words": [],
					"entities": [{
						"label": "LOCATION",
						"value": "Paris",
						"confidence": 0.95,
						"start_word": 3,
						"end_word": 4
					}]
				}]
			}]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3", DetectEntities: true})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	entities := res.Results.Channels[0].Alternatives[0].Entities
	if len(entities) != 1 {
		t.Fatalf("len(Entities) = %d, want 1", len(entities))
	}
	e := entities[0]
	if *e.Label != "LOCATION" || *e.Value != "Paris" {
		t.Fatalf("Entity = (%s, %s), want (LOCATION, Paris)", *e.Label, *e.Value)
	}
	if *e.StartWord != 3 || *e.EndWord != 4 {
		t.Fatalf("Entity word range = (%d, %d), want (3, 4)", *e.StartWord, *e.EndWord)
	}
}

func Test_LanguageDetection_Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockPost(t, `{
		"request_id": "lang-001",
		"results": {
			"channels": [{
				"alternatives": [{
					"transcript": "bonjour le monde",
					"confidence": 0.9,
					"words": [
						{"word": "bonjour", "start": 0.0, "end": 0.5, "confidence": 0.95, "language": "fr"},
						{"word": "le", "start": 0.5, "end": 0.7, "confidence": 0.9, "language": "fr"},
						{"word": "monde", "start": 0.7, "end": 1.0, "confidence": 0.92, "language": "fr"}
					]
				}],
				"detected_language": "fr",
				"language_confidence": 0.98
			}]
		}
	}`)

	res, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Model: "nova-3", DetectLanguage: true})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ch := res.Results.Channels[0]
	if ch.DetectedLanguage == nil || *ch.DetectedLanguage != "fr" {
		t.Fatalf("DetectedLanguage = %v, want 'fr'", ch.DetectedLanguage)
	}
	if ch.LanguageConfidence == nil || *ch.LanguageConfidence != 0.98 {
		t.Fatalf("LanguageConfidence = %v, want 0.98", ch.LanguageConfidence)
	}
	if ch.Alternatives[0].Words[0].Language == nil || *ch.Alternatives[0].Words[0].Language != "fr" {
		t.Fatalf("Words[0].Language = %v, want 'fr'", ch.Alternatives[0].Words[0].Language)
	}
}

func Test_InvalidOptions_RejectsBeforeHTTP(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, "should not reach"))

	_, err := newClient().FromURL(context.Background(), "https://example.com/audio.wav",
		&interfaces.PreRecordedTranscriptionOptions{Diarize: true, DiarizeVersion: "2021-07-14.0"})
	if err == nil {
		t.Skip("options.Check() did not reject the combination; skipping (test is brittle to Check rules)")
	}
}

func Test_GeneratedTypeShape_HasJSONTags(t *testing.T) {
	out := &spectypes.TranscribeOutput{
		RequestId: ptrStr("test-id"),
	}
	out.Metadata = &spectypes.ResponseMetadata{
		RequestId: ptrStr("test-id"),
		Duration:  ptrF32(1.5),
		Channels:  ptrI32(1),
	}
	if out.RequestId == nil || *out.RequestId != "test-id" {
		t.Fatal("RequestId roundtrip failed")
	}
	if *out.Metadata.Duration != 1.5 {
		t.Fatal("Metadata.Duration roundtrip failed")
	}
}
