// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// SageMaker transport-selection safety-net. Real InvokeEndpoint calls need live
// AWS + a deployed Deepgram model, so this verifies only that the option wires
// the SageMaker binding (cloud->self-hosted switch at construction). The
// transport's request/response mapping is covered by the generated binding.
package speakv1

import "testing"

func TestSageMakerTransportSelected(t *testing.T) {
	c := New(WithSageMakerTransport(nil, "deepgram-speak-endpoint"))
	if _, ok := c.transport.(sageMakerBinding); !ok {
		t.Fatalf("WithSageMakerTransport should select the SageMaker binding, got %T", c.transport)
	}
}
