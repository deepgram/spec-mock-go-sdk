// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"sync"
	"time"

	msginterface "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket/interfaces"
	common "github.com/deepgram/spec-mock-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces"
)

// internal structs
type controlMessage struct {
	Type string `json:"type"`
}

// Client is an alias for WSCallback
// Deprecated: use WSCallback instead
type Client = WSCallback

// WSCallback is a struct representing the websocket client connection using callbacks
type WSCallback struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	tOptions *interfaces.LiveTranscriptionOptions

	callback msginterface.LiveMessageCallback
	router   *commoninterfaces.Router

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
}

// WSChannel is a struct representing the websocket client connection using channels
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	tOptions *interfaces.LiveTranscriptionOptions

	chans  []*msginterface.LiveMessageChan
	router *commoninterfaces.Router

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
}

/*
Using Channels
*/
// DefaultChanHandler is a default channel handler for live transcription
// Simply prints the transcript to stdout
type DefaultChanHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool

	openChan          chan *msginterface.OpenResponse
	messageChan       chan *msginterface.MessageResponse
	metadataChan      chan *msginterface.MetadataResponse
	speechStartedChan chan *msginterface.SpeechStartedResponse
	utteranceEndChan  chan *msginterface.UtteranceEndResponse
	closeChan         chan *msginterface.CloseResponse
	errorChan         chan *msginterface.ErrorResponse
	unhandledChan     chan *[]byte
}

// ChanRouter routes events to channels
type ChanRouter struct {
	debugWebsocket bool

	openChan          []*chan *msginterface.OpenResponse
	messageChan       []*chan *msginterface.MessageResponse
	metadataChan      []*chan *msginterface.MetadataResponse
	speechStartedChan []*chan *msginterface.SpeechStartedResponse
	utteranceEndChan  []*chan *msginterface.UtteranceEndResponse
	closeChan         []*chan *msginterface.CloseResponse
	errorChan         []*chan *msginterface.ErrorResponse
	unhandledChan     []*chan *[]byte
}

/*
Using Callbacks
*/
// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultCallbackHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool
}

// CallbackRouter routes events to callbacks
type CallbackRouter struct {
	debugWebsocket bool
	callback       msginterface.LiveMessageCallback
}

// MessageRouter is the interface for routing messages
// Deprecated: Use CallbackRouter instead
type MessageRouter = CallbackRouter
