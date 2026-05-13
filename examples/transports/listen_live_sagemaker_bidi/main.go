// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Demonstrates the typed surface of the SageMaker BiDi streaming
// transport for listen-live. The transport package compiles and exports
// the Open*/Send/Recv/Close shape a real client expects, but the bodies
// currently return a "pending aws-sdk-go-v2 support" error because
// upstream hasn't released InvokeEndpointWithBidirectionalStream yet.
//
// When the SDK ships support, spec-codegen-go fills in the bodies and
// this example becomes runnable. Customer code can already depend on
// the package today.
package main

import (
	"github.com/deepgram/spec-mock-go-sdk/api/document"
	"context"
	"fmt"

	sm "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func main() {
	ctx := context.Background()
	endpoint := "your-deepgram-listen-live-endpoint"

	stream, err := sm.OpenStream[spectypes.ClientStream, spectypes.ServerStream](
		ctx, endpoint,
		spectypes.MarshalClientStream,
		spectypes.UnmarshalServerStream)
	if err != nil {
		fmt.Printf("sm.OpenStream: %v\n", err)
		return
	}
	defer stream.Close()

	configure := &spectypes.ClientStreamMemberConfigure{
		Value: spectypes.Configure{Features: map[string]document.Interface{}},
	}
	if err := stream.Send(configure); err != nil {
		fmt.Printf("Send Configure: %v\n", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Printf("Recv: %v\n", err)
			return
		}
		fmt.Printf("got: %T\n", msg)
	}
}
