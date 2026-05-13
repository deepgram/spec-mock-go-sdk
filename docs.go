// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Package sdk is the top-level marker for spec-mock-go-sdk, the Go
// consumer of the Deepgram Smithy spec pipeline. The package itself
// holds no exported symbols; its only job is the blank-import block
// below, which forces every client and api/-side package to compile
// as part of `go build ./...`. Add a new package here when its
// compilation should gate the build.
//
// Repo layout:
//   - api/   machine-generated wire types and transports
//   - pkg/   idiomatic Go facade (the part customers consume)
//
// See AGENTS.md and .agents/skills/ for the maintainer-facing
// conventions. Not for customer use; see deepgram/deepgram-go-sdk
// for the official Go SDK.
package sdk

import (
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/analyze"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/manage"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/client/speak"

	_ "github.com/deepgram/spec-mock-go-sdk/pkg/api/analyze/v1"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/api/manage/v1"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/api/speak/v1/rest"
	_ "github.com/deepgram/spec-mock-go-sdk/pkg/api/speak/v1/websocket"
)
