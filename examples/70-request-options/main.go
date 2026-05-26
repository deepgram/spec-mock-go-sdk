// Example: Additional request options
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
	client := prerecorded.New(prerecorded.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
	opts := &prerecorded.PreRecordedTranscriptionOptions{
		Model: "nova-3",
		AdditionalQueryParams: url.Values{
			"model":        []string{"nova-2-meeting"},
			"custom_param": []string{"custom-value"},
		},
	}

	resp, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.RequestID)
}
