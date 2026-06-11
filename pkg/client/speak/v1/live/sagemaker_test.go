// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// SageMaker bidi transport-selection safety-net (see the batch equivalent for
// why this only checks binding selection, not a live InvokeEndpoint stream).
package livev1

import "testing"

func TestSageMakerTransportSelected(t *testing.T) {
	c := New(WithSageMakerBidiTransport(nil, "deepgram-speak-live-endpoint"))
	if _, ok := c.transport.(sageMakerBidiBinding); !ok {
		t.Fatalf("WithSageMakerBidiTransport should select the SageMaker binding, got %T", c.transport)
	}
}
