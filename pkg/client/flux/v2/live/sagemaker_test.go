// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// SageMaker bidi transport-selection safety-net for Flux (Listen v2).
package livev2

import "testing"

func TestSageMakerTransportSelected(t *testing.T) {
	c := New(WithSageMakerBidiTransport(nil, "deepgram-flux-endpoint"))
	if _, ok := c.transport.(sageMakerBidiBinding); !ok {
		t.Fatalf("WithSageMakerBidiTransport should select the SageMaker binding, got %T", c.transport)
	}
}
