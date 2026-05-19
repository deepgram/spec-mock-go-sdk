// Example: Error Handling
//
// Two distinct error patterns for two transports:
//
//   REST: typed *HTTPError from api/transport/http. Use errors.As to
//         extract the type, then discriminate on StatusCode for
//         status-specific handling (retry on 429, re-auth on 401, etc.).
//
//   WebSocket: server-emitted errors arrive as *wsv1.ErrorEvent through
//         stream.Recv(), NOT as exceptions. Facade-level send-side
//         failures (frame too large, send timeout) use sentinel-style
//         exported errors matched via errors.Is.
//
// Plus context cancellation, which applies to both transports.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
	wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func main() {
	example1HTTPErrorExtraction()
	example2DiscriminateByStatusCode()
	example3TypedErrorViaErrorsAs()
	example4NonHTTPNetworkErrors()
	example5WebSocketErrors()
}

func example1HTTPErrorExtraction() {
	fmt.Println("Example 1: Extract typed *HTTPError via errors.As")
	client := restv1.NewWithDefaults()

	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	if errors.As(err, &httpErr) {
		fmt.Println("  HTTP error:")
		fmt.Printf("    StatusCode: %d\n", httpErr.StatusCode)
		fmt.Printf("    Body:       %s\n", string(httpErr.Body))
		fmt.Printf("    Method:     %s\n", httpErr.Method)
		fmt.Printf("    URL:        %s\n", httpErr.URL)
		fmt.Printf("    Content-Type header: %s\n", httpErr.Headers.Get("Content-Type"))
	} else if err != nil {
		fmt.Printf("  Non-HTTP error (network/transport): %v\n", err)
	}
}

func example2DiscriminateByStatusCode() {
	fmt.Println("\nExample 2: Status-specific handling via StatusCode switch")
	client := restv1.NewWithDefaults()

	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		fmt.Printf("  Not an HTTP error: %v\n", err)
		return
	}
	switch httpErr.StatusCode {
	case http.StatusBadRequest:
		fmt.Println("  400 Bad Request — caller-side problem; fix the request and retry.")
	case http.StatusUnauthorized:
		fmt.Println("  401 Unauthorized — credentials missing or invalid; re-auth and retry.")
	case http.StatusForbidden:
		fmt.Println("  403 Forbidden — credentials valid but lack scope for this operation.")
	case http.StatusTooManyRequests:
		fmt.Println("  429 Too Many Requests — back off using the Retry-After header and retry.")
		if retryAfter := httpErr.Headers.Get("Retry-After"); retryAfter != "" {
			fmt.Printf("    Retry-After: %s\n", retryAfter)
		}
	default:
		switch {
		case httpErr.StatusCode >= 500:
			fmt.Printf("  %d server error — retry with exponential backoff.\n", httpErr.StatusCode)
		default:
			fmt.Printf("  %d %s\n", httpErr.StatusCode, string(httpErr.Body))
		}
	}
}

func example3TypedErrorViaErrorsAs() {
	fmt.Println("\nExample 3: Typed error via errors.As on HTTPError.Typed")
	client := restv1.NewWithDefaults()

	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	if !errors.As(err, &httpErr) {
		fmt.Printf("  Not an HTTP error: %v\n", err)
		return
	}
	if httpErr.Typed == nil {
		fmt.Printf("  Typed nil — status %d is not in Transcribe's declared errors, or body was undecodable.\n", httpErr.StatusCode)
		return
	}

	var queryErr *spectypes.InvalidQueryParameterError
	if errors.As(httpErr.Typed, &queryErr) {
		fmt.Println("  Typed extracted as *spectypes.InvalidQueryParameterError:")
		if queryErr.ErrCode != nil {
			fmt.Printf("    ErrCode:   %s\n", *queryErr.ErrCode)
		}
		if queryErr.ErrMsg != nil {
			fmt.Printf("    ErrMsg:    %s\n", *queryErr.ErrMsg)
		}
		if queryErr.RequestId != nil {
			fmt.Printf("    RequestId: %s\n", *queryErr.RequestId)
		}
		return
	}

	var rateLimited *spectypes.RateLimitedError
	if errors.As(httpErr.Typed, &rateLimited) {
		fmt.Println("  Typed extracted as *spectypes.RateLimitedError")
		if rateLimited.RetryAfter != nil {
			fmt.Printf("    Retry-After: %s seconds\n", *rateLimited.RetryAfter)
		}
		return
	}

	fmt.Printf("  Typed is set but did not match expected variants: %T\n", httpErr.Typed)
}

func example4NonHTTPNetworkErrors() {
	fmt.Println("\nExample 4: Non-HTTP transport failures (DNS, connection refused, context)")
	client := restv1.NewWithDefaults()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.FromURL(
		ctx,
		"https://dpgr.am/spacewalk.wav",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)

	var httpErr *httptransport.HTTPError
	switch {
	case errors.As(err, &httpErr):
		fmt.Printf("  Unexpected HTTPError on a context-deadline test: %d\n", httpErr.StatusCode)
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("  Context deadline exceeded — request never completed in time.")
	case errors.Is(err, context.Canceled):
		fmt.Println("  Context cancelled by caller.")
	case err != nil:
		fmt.Printf("  Transport error (DNS / TCP / TLS): %v\n", err)
	default:
		fmt.Println("  Request succeeded before the deadline.")
	}
}

func example5WebSocketErrors() {
	fmt.Println("\nExample 5: WebSocket error handling")
	fmt.Println("")
	fmt.Println("  Server-emitted errors arrive through stream.Recv() as a")
	fmt.Println("  *wsv1.ErrorEvent variant of the Event union. Handle them in")
	fmt.Println("  the Recv loop's type-switch:")
	fmt.Println("")
	fmt.Println(`    for {`)
	fmt.Println(`        event, err := stream.Recv()`)
	fmt.Println(`        if errors.Is(err, io.EOF) { return }`)
	fmt.Println(`        if err != nil { return err }`)
	fmt.Println(`        switch e := event.(type) {`)
	fmt.Println(`        case *wsv1.ErrorEvent:`)
	fmt.Printf("            log.Printf(\"deepgram error: %%s\", e.Description)\n")
	fmt.Println(`            return`)
	fmt.Println(`        case *wsv1.ResultsEvent:`)
	fmt.Println(`            // handle transcript ...`)
	fmt.Println(`        }`)
	fmt.Println(`    }`)
	fmt.Println("")
	fmt.Println("  Facade-level send-side validation errors are sentinel-style,")
	fmt.Println("  matched via errors.Is. They surface BEFORE bytes hit the wire:")
	fmt.Println("")

	simulated := fmt.Errorf("%w: 200-byte chunk exceeds 100-byte cap", wsv1.ErrFrameTooLarge)
	if errors.Is(simulated, wsv1.ErrFrameTooLarge) {
		fmt.Println("    if errors.Is(sendErr, wsv1.ErrFrameTooLarge) { ... }")
	}

	timedOut := fmt.Errorf("%w: write blocked > 5s", wsv1.ErrSendTimeout)
	if errors.Is(timedOut, wsv1.ErrSendTimeout) {
		fmt.Println("    if errors.Is(sendErr, wsv1.ErrSendTimeout) { ... }")
	}
}
