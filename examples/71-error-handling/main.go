// Example: Error handling
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	httptransport "github.com/deepgram/spec-mock-go-sdk/api/transport/http"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	_, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model: "not-a-real-model",
	})
	if err == nil {
		return
	}

	var httpErr *httptransport.HTTPError
	if errors.As(err, &httpErr) {
		fmt.Println("status", httpErr.StatusCode)

		var queryErr *spectypes.InvalidQueryParameterError
		if errors.As(httpErr.Typed, &queryErr) {
			fmt.Println("query error", value(queryErr.ErrCode))
		}
		return
	}

	log.Fatal(err)
}

func value(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
