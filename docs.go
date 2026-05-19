// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Package sdk is the top-level marker for spec-mock-go-sdk, the Go
// consumer of the Deepgram Smithy spec pipeline. The package itself
// holds no exported symbols; its only job is the blank-import block
// below, which forces every package the build is supposed to gate
// on to compile as part of `go build ./...`.
//
// As products migrate into the spec pipeline they each get their own
// pkg/client/{product}/{ver}/{transport}/ package; add them here.
package sdk

import (
	_ "github.com/deepgram/spec-mock-go-sdk/api/document"
	_ "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	_ "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
	_ "github.com/deepgram/spec-mock-go-sdk/api/transport/webrtc"
	_ "github.com/deepgram/spec-mock-go-sdk/api/transport/websocket"
	_ "github.com/deepgram/spec-mock-go-sdk/api/types"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)
