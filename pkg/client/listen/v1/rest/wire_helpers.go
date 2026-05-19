// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT
//
// Wire-test helpers shared between the hand-written wire_test.go and the
// codegen-emitted wire_test_generated.go. These helpers live in a normal
// (non-_test.go) source file so the generated file — which uses the
// suffix _generated.go rather than _test.go — can resolve them at
// compile time. They are still only useful from test contexts (they take
// *testing.T), but they're symbol-visible to the whole package.
//
// If a future codegen regen renames wire_test_generated.go to
// wire_generated_test.go, these helpers can move back into wire_test.go;
// nothing else in the facade calls them.

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
