// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package wsv1

// Event is the union of customer-facing events emitted by a live
// transcription Stream. Customer code receives values implementing
// Event from Stream.Recv and type-switches on the concrete type to
// handle each kind:
//
//	for {
//	    event, err := stream.Recv()
//	    if err != nil { return err }
//	    switch e := event.(type) {
//	    case *ResultsEvent:
//	        fmt.Println(e.Channel.Alternatives[0].Transcript)
//	    case *MetadataEvent:
//	        // session ended
//	    case *ErrorEvent:
//	        // mid-stream error
//	    }
//	}
//
// Each Event variant mirrors one ServerStream member in the api/
// wire surface. New variants added by future spec evolutions get
// new concrete types here; existing customer type-switches keep
// compiling and just don't match the new variant.
type Event interface {
	isEvent()
}

// ResultsEvent carries an interim or final transcription result.
// Mirrors spectypes.StreamingResponse.
type ResultsEvent struct {
	Channel      *Channel
	ChannelIndex []int
	Duration     float64
	FromFinalize bool
	IsFinal      bool
	Metadata     *InterimMetadata
	Start        float64
	Entities     []Entity
	SpeechFinal  bool
	Tag          string
}

func (*ResultsEvent) isEvent() {}

// MetadataEvent carries per-session metadata. Always emitted at
// session end. Mirrors spectypes.WsMetadata.
type MetadataEvent struct {
	RequestID string
	Channels  int
	Created   string
	Duration  float64
	Extra     map[string]string
	ModelInfo map[string]ModelInfo
	Sha256    string
	Warnings  []Warning
}

func (*MetadataEvent) isEvent() {}

// SpeechStartedEvent is the speech-detected signal, emitted only
// when VadEvents=true. Mirrors spectypes.SpeechStarted.
type SpeechStartedEvent struct {
	Channel   []int
	Timestamp float64
}

func (*SpeechStartedEvent) isEvent() {}

// UtteranceEndEvent is the VAD-driven utterance boundary signal.
// Mirrors spectypes.UtteranceEnd.
type UtteranceEndEvent struct {
	Channel     []int
	LastWordEnd float64
}

func (*UtteranceEndEvent) isEvent() {}

// ErrorEvent carries a mid-stream or pre-close error. Mirrors
// spectypes.WsError.
type ErrorEvent struct {
	Variant     string
	Code        string
	Description string
	Message     string
}

func (*ErrorEvent) isEvent() {}

// SyncEvent is an echo of a client Sync message; used to align
// test pipelines. Mirrors spectypes.ServerSync.
type SyncEvent struct {
	ID int64
}

func (*SyncEvent) isEvent() {}

// Channel is the customer-facing single-channel transcript shape
// inside a ResultsEvent. Multichannel streams emit one ResultsEvent
// per channel, each carrying a Channel selected by ChannelIndex.
type Channel struct {
	Search             []Search
	Alternatives       []Alternative
	DetectedLanguage   string
	LanguageConfidence float64
}

type Alternative struct {
	Transcript string
	Confidence float64
	Words      []Word
	Languages  []string
}

type Word struct {
	Word              string
	Start             float64
	End               float64
	Confidence        float64
	Speaker           *int
	SpeakerConfidence *float64
	PunctuatedWord    string
	Language          string
}

type Search struct {
	Query string
	Hits  []Hit
}

type Hit struct {
	Confidence float64
	Start      float64
	End        float64
	Snippet    string
}

type Entity struct {
	Label      string
	Value      string
	Confidence float64
	StartWord  int
	EndWord    int
}

type InterimMetadata struct {
	RequestID string
	ModelUUID string
	ModelInfo *ModelInfo
	Extra     map[string]string
}

type ModelInfo struct {
	Name    string
	Version string
	Arch    string
}

type Warning struct {
	Parameter string
	Type      string
	Message   string
}
