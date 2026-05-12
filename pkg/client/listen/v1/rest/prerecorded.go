// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package restv1

import (
	"context"
	"io"

	klog "k8s.io/klog/v2"

	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	interfaces "github.com/deepgram/spec-mock-go-sdk/pkg/client/interfaces"
)

// FromFile transcribes a prerecorded audio file from a local path.
func (c *Client) FromFile(ctx context.Context, file string, options *interfaces.PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	return c.sendAudio(ctx, func(ctx context.Context, opts *interfaces.PreRecordedTranscriptionOptions, resp interface{}) error {
		return c.DoFile(ctx, file, opts, resp)
	}, options)
}

// FromStream transcribes a prerecorded audio file from a reader.
func (c *Client) FromStream(ctx context.Context, src io.Reader, options *interfaces.PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	return c.sendAudio(ctx, func(ctx context.Context, opts *interfaces.PreRecordedTranscriptionOptions, resp interface{}) error {
		return c.DoStream(ctx, src, opts, resp)
	}, options)
}

// FromURL transcribes a prerecorded audio file fetched from a remote URL.
func (c *Client) FromURL(ctx context.Context, url string, options *interfaces.PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	return c.sendAudio(ctx, func(ctx context.Context, opts *interfaces.PreRecordedTranscriptionOptions, resp interface{}) error {
		return c.DoURL(ctx, url, opts, resp)
	}, options)
}

type sendFunc func(ctx context.Context, options *interfaces.PreRecordedTranscriptionOptions, resp interface{}) error

// sendAudio runs the request through the generated spectypes.TranscribeOutput
// shape (so json.Unmarshal exercises the spec-generated JSON tags) and then
// converts to the idiomatic *PreRecordedResponse value-typed facade before
// handing back to the caller. The two-step shape is intentional: it puts the
// generated wire types on the parse hot path while keeping the
// customer-facing API stable across api/types regenerations.
func (c *Client) sendAudio(ctx context.Context, sender sendFunc, options *interfaces.PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	klog.V(6).Infof("prerecorded.sendAudio ENTER\n")

	if err := options.Check(); err != nil {
		klog.V(1).Infof("PreRecordedTranscriptionOptions.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("prerecorded.sendAudio LEAVE\n")
		return nil, err
	}

	var wire spectypes.TranscribeOutput

	if err := sender(ctx, options, &wire); err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
		}
		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("prerecorded.sendAudio LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("sendAudio Succeeded\n")
	klog.V(6).Infof("prerecorded.sendAudio LEAVE\n")
	return convertTranscribeOutput(&wire), nil
}
