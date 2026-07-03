// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// SAGEMAKER-001 safety-net: the batch SageMaker binding has no URL query
// string, so typed query options must ride through CustomAttributes as the
// request target "<path>?<query>". Drives the real binding against a fake
// SageMaker runtime endpoint and asserts the header the AWS SDK emits.
package speakv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
)

func TestSageMakerBatchForwardsTypedQuery(t *testing.T) {
	var gotAttrs, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAttrs = r.Header.Get("X-Amzn-SageMaker-Custom-Attributes")
		gotContentType = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "audio/mpeg")
		_, _ = w.Write([]byte{0x01, 0x02, 0x03})
	}))
	defer srv.Close()

	cfg := aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("AKID", "SECRET", ""),
	}
	smClient := sagemakerruntime.NewFromConfig(cfg, func(o *sagemakerruntime.Options) {
		o.BaseEndpoint = aws.String(srv.URL)
	})

	c := New(WithSageMakerTransport(smClient, "deepgram-speak-endpoint"))
	out, err := c.FromText(context.Background(), "hello",
		&SpeakOptions{Model: "aura-2-asteria-en", Encoding: "mp3"})
	if err != nil {
		t.Fatalf("FromText via SageMaker: %v", err)
	}

	// url.Values.Encode sorts keys, so the order is deterministic.
	if want := "v1/speak?encoding=mp3&model=aura-2-asteria-en"; gotAttrs != want {
		t.Fatalf("CustomAttributes = %q, want %q", gotAttrs, want)
	}
	if gotContentType != "text/plain" {
		t.Fatalf("Content-Type = %q, want text/plain", gotContentType)
	}
	if string(out.Audio) != string([]byte{0x01, 0x02, 0x03}) {
		t.Fatalf("audio body = %v, want the response bytes", out.Audio)
	}
}

func TestSageMakerBatchNoOptions_PathOnly(t *testing.T) {
	var gotAttrs string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAttrs = r.Header.Get("X-Amzn-SageMaker-Custom-Attributes")
		w.Header().Set("Content-Type", "audio/mpeg")
		_, _ = w.Write([]byte{0x00})
	}))
	defer srv.Close()

	cfg := aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("AKID", "SECRET", ""),
	}
	smClient := sagemakerruntime.NewFromConfig(cfg, func(o *sagemakerruntime.Options) {
		o.BaseEndpoint = aws.String(srv.URL)
	})

	c := New(WithSageMakerTransport(smClient, "deepgram-speak-endpoint"))
	if _, err := c.FromText(context.Background(), "hi", nil); err != nil {
		t.Fatalf("FromText via SageMaker: %v", err)
	}
	if gotAttrs != "v1/speak" {
		t.Fatalf("CustomAttributes = %q, want bare path v1/speak", gotAttrs)
	}
}
