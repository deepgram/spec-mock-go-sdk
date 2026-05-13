// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	msginterface "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket/interfaces"
)

// NewWithDefault creates a CallbackRouter with the default callback handler
// Deprecated: Use NewCallbackWithDefault instead
func NewWithDefault() *CallbackRouter {
	return NewCallbackWithDefault()
}

// NewCallbackWithDefault creates a CallbackRouter with the default callback handler
func NewCallbackWithDefault() *CallbackRouter {
	return NewCallbackRouter(NewDefaultCallbackHandler())
}

// NewCallbackRouter creates a CallbackRouter with a user-defined callback
func NewCallbackRouter(callback msginterface.LiveMessageCallback) *CallbackRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}
	return &CallbackRouter{
		callback:       callback,
		debugWebsocket: strings.EqualFold(strings.ToLower(debugStr), "true"),
	}
}

// Open sends an OpenResponse message to the callback
func (r *CallbackRouter) Open(or *msginterface.OpenResponse) error {
	return r.callback.Open(or)
}

// Close sends an CloseResponse message to the callback
func (r *CallbackRouter) Close(or *msginterface.CloseResponse) error {
	return r.callback.Close(or)
}

// Error sends an ErrorResponse message to the callback
func (r *CallbackRouter) Error(er *msginterface.ErrorResponse) error {
	return r.callback.Error(er)
}

// processMessage generalizes the handling of all message types
func (r *CallbackRouter) processGeneric(msgType string, byMsg []byte, action func(data *interface{}) error, data interface{}) error {
	klog.V(6).Infof("router.%s ENTER\n", msgType)

	r.printDebugMessages(5, msgType, byMsg)

	var err error
	if err = action(&data); err != nil {
		klog.V(1).Infof("callback.%s failed. Err: %v\n", msgType, err)
	} else {
		klog.V(5).Infof("callback.%s succeeded\n", msgType)
	}
	klog.V(6).Infof("router.%s LEAVE\n", msgType)

	return err
}

// Message handles platform messages and routes them appropriately based on the MessageType
func (r *CallbackRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("router.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	// Generated dispatcher: see chan_router.go's Message() for the
	// architectural rationale. Same approach mirrored here for the
	// callback-based router.
	msg, err := spectypes.UnmarshalServerStream(byMsg)
	if err != nil {
		klog.V(1).Infof("UnmarshalServerStream failed. Err: %v\n", err)
		klog.V(6).Infof("router.Message LEAVE\n")
		return r.UnhandledMessage(byMsg)
	}

	switch m := msg.(type) {
	case *spectypes.ServerStreamMemberResults:
		err = r.callback.Message(msginterface.FromStreamingResponse(&m.Value))
	case *spectypes.ServerStreamMemberMetadata:
		err = r.callback.Metadata(msginterface.FromWsMetadata(&m.Value))
	case *spectypes.ServerStreamMemberSpeechStarted:
		err = r.callback.SpeechStarted(msginterface.FromSpeechStarted(&m.Value))
	case *spectypes.ServerStreamMemberUtteranceEnd:
		err = r.callback.UtteranceEnd(msginterface.FromUtteranceEnd(&m.Value))
	case *spectypes.ServerStreamMemberError:
		err = r.callback.Error(wsErrorToSDKError(&m.Value))
	case *spectypes.ServerStreamMemberSync:
		// Sync messages are test-pipeline alignment markers, not
		// part of the public surface. Ignore.
	default:
		err = r.UnhandledMessage(byMsg)
	}

	if err != nil {
		klog.V(5).Infof("router.Message dispatch failed: %v\n", err)
	}
	klog.V(6).Infof("router.Message LEAVE\n")
	return err
}

// Binary handles platform messages and routes them appropriately based on the MessageType
func (r *CallbackRouter) Binary(byMsg []byte) error {
	// No implementation needed on STT
	return nil
}

// UnhandledMessage logs and handles any unexpected message types
func (r *CallbackRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("router.UnhandledMessage ENTER\n")
	r.printDebugMessages(3, "UnhandledMessage", byMsg)
	klog.V(1).Infof("Unknown Event was received\n")
	klog.V(6).Infof("router.UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

// printDebugMessages formats and logs debugging messages
func (r *CallbackRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJSON, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Format failed. Err: %v\n", err)
		return
	}
	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJSON)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
