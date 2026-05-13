// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1_test

import (
	"fmt"

	rest "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

// ExampleNewWithDefaults constructs a Listen v1 REST client using the
// env-backed credential layer. Set DEEPGRAM_API_KEY (or
// DEEPGRAM_ACCESS_TOKEN) in the environment and the client picks them
// up at construction time.
func ExampleNewWithDefaults() {
	client := rest.NewWithDefaults()
	_ = client
	fmt.Println("constructed listen v1 rest client")
	// Output: constructed listen v1 rest client
}

// ExamplePreRecordedResponse demonstrates the customer-facing shape
// of a prerecorded transcription response. Primitive fields are
// value-typed and zero-valued when absent, so reads need no nil
// check. Nested struct fields like Metadata and Results stay
// pointer-typed and may be nil — guard them before dereferencing.
func ExamplePreRecordedResponse() {
	resp := &rest.PreRecordedResponse{
		RequestID: "req-abc123",
	}
	fmt.Println("RequestID:", resp.RequestID)
	if resp.Metadata == nil {
		fmt.Println("no metadata on this response")
	}
	// Output:
	// RequestID: req-abc123
	// no metadata on this response
}

// ExampleMetadata shows the Metadata sub-shape returned alongside a
// PreRecordedResponse. Duration is a float64 (widened from the
// generated float32 by the facade) and Channels is a plain int.
func ExampleMetadata() {
	m := &rest.Metadata{
		RequestID: "req-abc123",
		Duration:  12.345,
		Channels:  2,
	}
	fmt.Printf("request=%s duration=%.3fs channels=%d\n",
		m.RequestID, m.Duration, m.Channels)
	// Output: request=req-abc123 duration=12.345s channels=2
}

// ExampleIsURL distinguishes a remote URL from a local file path.
// PreRecordedClient.FromURL uses this to decide whether to fetch the
// audio remotely or open it from disk.
func ExampleIsURL() {
	fmt.Println(rest.IsURL("https://dpgr.am/spacewalk.wav"))
	fmt.Println(rest.IsURL("/local/audio.wav"))
	// Output:
	// true
	// false
}
