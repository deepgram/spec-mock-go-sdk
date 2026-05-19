// Example: Error Handling
//
// Shows the common patterns for handling errors from the SDK:
// typed HTTPError.Typed via errors.As, the raw HTTPError fields when
// Typed is nil, context cancellation, and the WebSocket Recv-loop
// type-switch plus send-side sentinels.

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
	wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func main() {
	typedError()
	rawHTTPErrorFallback()
	nonHTTPError()
	webSocketError()
}

func typedError() {
	client := restv1.NewWithDefaults()

	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		return
	}

	var queryErr *spectypes.InvalidQueryParameterError
	if errors.As(httpErr.Typed, &queryErr) {
		fmt.Println("InvalidQueryParameterError:")
		printField("err_code", queryErr.ErrCode)
		printField("err_msg", queryErr.ErrMsg)
		printField("request_id", queryErr.RequestId)
		return
	}

	var rateLimited *spectypes.RateLimitedError
	if errors.As(httpErr.Typed, &rateLimited) {
		fmt.Println("RateLimitedError:")
		printField("Retry-After", rateLimited.RetryAfter)
		return
	}
}

func rawHTTPErrorFallback() {
	client := restv1.NewWithDefaults()

	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	if errors.As(err, &httpErr) && httpErr.Typed == nil {
		fmt.Printf("Untyped %d response from %s: %s\n",
			httpErr.StatusCode, httpErr.URL, string(httpErr.Body))
	}
}

func nonHTTPError() {
	client := restv1.NewWithDefaults()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.FromURL(
		ctx,
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("deadline exceeded")
	case errors.Is(err, context.Canceled):
		fmt.Println("cancelled")
	case err != nil:
		fmt.Printf("transport error: %v\n", err)
	}
}

func webSocketError() {
	frameTooBig := fmt.Errorf("%w: chunk too big", wsv1.ErrFrameTooLarge)
	if errors.Is(frameTooBig, wsv1.ErrFrameTooLarge) {
		fmt.Println("errors.Is(sendErr, wsv1.ErrFrameTooLarge)")
	}

	timedOut := fmt.Errorf("%w: write blocked", wsv1.ErrSendTimeout)
	if errors.Is(timedOut, wsv1.ErrSendTimeout) {
		fmt.Println("errors.Is(sendErr, wsv1.ErrSendTimeout)")
	}
}

func printField(label string, p *string) {
	if p == nil {
		return
	}
	fmt.Printf("  %s: %s\n", label, *p)
}
