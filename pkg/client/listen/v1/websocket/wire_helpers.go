// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT
//
// Wire-test helpers for the live transcription facade. Mirror of
// pkg/client/listen/v1/rest/wire_helpers.go but against StreamInput.

package wsv1

import (
	"reflect"
	"testing"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func requireWired(t *testing.T, in *spectypes.StreamInput, fieldName string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		t.Fatalf("spectypes.StreamInput has no field %q yet — spec needs to model it before the converter can wire it through.", fieldName)
	}
	if isZeroForWire(v) {
		t.Fatalf("spectypes.StreamInput.%s exists but optionsToStreamInput didn't wire it.", fieldName)
	}
}

func requireDropped(t *testing.T, in *spectypes.StreamInput, fieldName, reason string) {
	t.Helper()
	v := reflect.ValueOf(in).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		return
	}
	if !isZeroForWire(v) {
		t.Fatalf("spectypes.StreamInput.%s is documented as permanently dropped (%s) but the converter wired it anyway.", fieldName, reason)
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
