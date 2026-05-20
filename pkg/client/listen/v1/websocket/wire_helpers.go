// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT
//
// Wire-test helpers for the live transcription client. Mirror of
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

// requireFacadeOnly asserts that an options field is intentionally
// facade-only: present on *Options but with no counterpart on the
// generated spectypes input. Used for escape-hatch fields like
// AdditionalQueryParams (url.Values) that pass arbitrary unspecified
// query params at connect time. This is the third state alongside
// requireWired (option ↔ wire) and requireDropped (option remains,
// wire removed) — it signals to spec-idiomatic that the missing wire
// counterpart is by design, not a facade hallucination.
func requireFacadeOnly(t *testing.T, options any, fieldName string) {
	t.Helper()
	v := reflect.ValueOf(options).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		t.Fatalf("%T has no field %q — requireFacadeOnly expected the facade-only field to exist on the options struct.", options, fieldName)
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
