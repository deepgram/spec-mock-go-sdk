// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Demonstrates the SageMaker batch transport for the listen REST surface.
// Reads DG_SAGEMAKER_ENDPOINT for the endpoint name and the AWS environment
// (DEFAULT credentials, region) the way every aws-sdk-go-v2 caller does.
//
// Runs against a real Deepgram-on-SageMaker endpoint if you have one; the
// httptest-driven unit test in tests/unit_test/transport_test.go exercises
// the same code path without external dependencies.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"

	sm "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func main() {
	endpoint := os.Getenv("DG_SAGEMAKER_ENDPOINT")
	if endpoint == "" {
		log.Fatal("DG_SAGEMAKER_ENDPOINT not set")
	}

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}
	smr := sagemakerruntime.NewFromConfig(awsCfg)

	model := "nova-3"
	language := "en-US"
	input := &spectypes.TranscribeInput{
		Model:    &model,
		Language: &language,
	}

	out, err := sm.InvokeTranscribe(ctx, smr, endpoint, input)
	if err != nil {
		log.Fatalf("InvokeTranscribe: %v", err)
	}

	if out.RequestId != nil {
		fmt.Printf("Request ID: %s\n", *out.RequestId)
	}
	if out.Results != nil && len(out.Results.Channels) > 0 &&
		len(out.Results.Channels[0].Alternatives) > 0 &&
		out.Results.Channels[0].Alternatives[0].Transcript != nil {
		fmt.Printf("Transcript: %s\n", *out.Results.Channels[0].Alternatives[0].Transcript)
	}
}
