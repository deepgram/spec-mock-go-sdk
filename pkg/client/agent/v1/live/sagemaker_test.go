// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// SageMaker bidi transport-selection safety-net for the Voice Agent.
package livev1

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestSageMakerTransportSelected(t *testing.T) {
	c := New(WithSageMakerBidiTransport(aws.Config{}, "deepgram-agent-endpoint"))
	if _, ok := c.transport.(sageMakerBidiBinding); !ok {
		t.Fatalf("WithSageMakerBidiTransport should select the SageMaker binding, got %T", c.transport)
	}
}
