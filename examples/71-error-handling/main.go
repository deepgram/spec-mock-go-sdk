// Example: Error Handling
//
// Go errors are plain values; the idiomatic patterns are errors.Is
// (sentinel match), errors.As (type assertion), and context
// cancellation. This SDK follows that style — no custom Error class
// hierarchy.
//
// Where the SDK exports sentinel errors:
//   - wsv1.ErrFrameTooLarge: SendAudio chunk exceeds Config.MaxFrameSizeBytes
//   - wsv1.ErrSendTimeout:   SendAudio exceeded Config.SendTimeout
//
// Everything else surfaces as a wrapped error containing the HTTP
// status / Deepgram error body / underlying transport error.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
	wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

func main() {
	example1HandlingApiErrors()
	example2InvalidRequest()
	example3ContextCancellation()
	example4WebSocketSentinels()
}

func example1HandlingApiErrors() {
	fmt.Println("Example 1: Handling API errors (bad URL)")
	client := restv1.NewWithDefaults()
	_, err := client.FromURL(
		context.Background(),
		"https://invalid-url.example.com/audio.wav",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)
	if err != nil {
		fmt.Printf("  Got expected error: %v\n", err)
	}
}

func example2InvalidRequest() {
	fmt.Println("\nExample 2: Empty URL (server returns 400)")
	client := restv1.NewWithDefaults()
	_, err := client.FromURL(
		context.Background(),
		"",
		&restv1.PreRecordedTranscriptionOptions{Model: "nova-3"},
	)
	if err != nil {
		fmt.Printf("  Got expected error: %v\n", err)
	}
}

func example3ContextCancellation() {
	fmt.Println("\nExample 3: Context cancellation")
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
		fmt.Println("  Caught context.DeadlineExceeded — call cancelled by client.")
	case err != nil:
		fmt.Printf("  Other error: %v\n", err)
	default:
		fmt.Println("  Surprisingly fast — call completed before the 1ms deadline.")
	}
}

func example4WebSocketSentinels() {
	fmt.Println("\nExample 4: WebSocket sentinel errors")

	err := fmt.Errorf("%w: simulated 200-byte frame exceeds 100-byte limit",
		wsv1.ErrFrameTooLarge)

	if errors.Is(err, wsv1.ErrFrameTooLarge) {
		fmt.Println("  Caught wsv1.ErrFrameTooLarge — pattern matches via errors.Is.")
		fmt.Println("  In production, this error comes from stream.SendAudio when")
		fmt.Println("  Client.WithConfig(&Config{MaxFrameSizeBytes: N}) is set and")
		fmt.Println("  a chunk exceeds N bytes.")
	} else if err != nil {
		log.Printf("  Got unexpected error: %v", err)
	}
}
