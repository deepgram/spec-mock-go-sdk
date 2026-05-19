// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT
//
// Wire-test helpers for the prerecorded client. Used by wire_test.go.
// Mirror of pkg/client/listen/v1/websocket/wire_helpers.go but against
// TranscribeInput.

package restv1

import (
	"reflect"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func requireWired(t *testing.T, in *spectypes.TranscribeInput, fieldName string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		t.Fatalf("spectypes.TranscribeInput has no field %q yet — spec needs to model it before the converter can wire it through.", fieldName)
	}
	if isZeroForWire(v) {
		t.Fatalf("spectypes.TranscribeInput.%s exists but optionsToTranscribeInput didn't wire it.", fieldName)
	}
}

func requireDropped(t *testing.T, in *spectypes.TranscribeInput, fieldName, reason string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		return
	}
	if !isZeroForWire(v) {
		t.Fatalf("spectypes.TranscribeInput.%s is documented as permanently dropped (%s) but the converter wired it anyway.", fieldName, reason)
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
