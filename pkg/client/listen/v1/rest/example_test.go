// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1_test

import (
	"fmt"

	rest "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func ExampleNewWithDefaults() {
	client := rest.NewWithDefaults()
	_ = client
	fmt.Println("constructed listen v1 rest client")
	// Output: constructed listen v1 rest client
}

func ExamplePreRecordedTranscriptionOptions() {
	opts := &rest.PreRecordedTranscriptionOptions{
		Model:       "nova-3",
		Language:    "en-US",
		Punctuate:   true,
		SmartFormat: true,
		Utterances:  true,
	}
	_ = opts
	fmt.Println("ok")
	// Output: ok
}
