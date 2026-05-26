// Example: Prerecorded transcription with SageMaker
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
	awsClient := sagemakerruntime.New(sagemakerruntime.Options{Region: "us-east-1"})
	client := prerecorded.New(prerecorded.WithSageMakerTransport(
		awsClient,
		"endpoint-name",
		prerecorded.WithTargetVariant("variantA"),
	))

	resp, err := client.FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model: "nova-3",
		AdditionalQueryParams: url.Values{
			"request_id": []string{"demo-123"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.RequestID)
}
