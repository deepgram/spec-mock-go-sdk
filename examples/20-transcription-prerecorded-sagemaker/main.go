// Example: Prerecorded transcription with SageMaker
// SageMaker endpoints typically run with network isolation enabled, so the
// container has no outbound internet access and cannot fetch a remote audio
// URL itself. The client downloads the audio and streams the bytes through
// InvokeEndpoint.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	prerecorded "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/prerecorded"
)

func main() {
	endpoint := os.Getenv("SAGEMAKER_ENDPOINT_NAME")
	if endpoint == "" {
		log.Fatal("SAGEMAKER_ENDPOINT_NAME is required")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}
	awsClient := sagemakerruntime.NewFromConfig(cfg)

	client := prerecorded.New(prerecorded.WithSageMakerTransport(
		awsClient,
		endpoint,
	))

	audioURL := os.Getenv("AUDIO_URL")
	if audioURL == "" {
		audioURL = "https://dpgr.am/spacewalk.wav"
	}
	audioPath := os.Getenv("AUDIO_PATH")
	var src io.ReadCloser
	if audioPath != "" {
		f, err := os.Open(audioPath)
		if err != nil {
			log.Fatalf("open audio: %v", err)
		}
		src = f
	} else {
		resp, err := http.Get(audioURL)
		if err != nil {
			log.Fatalf("download audio: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("download audio: %s", resp.Status)
		}
		src = resp.Body
	}
	defer src.Close()

	out, err := client.FromStream(ctx, src, "audio/wav", &prerecorded.PreRecordedTranscriptionOptions{
		Model:       "nova-3",
		Punctuate:   true,
		SmartFormat: true,
		AdditionalQueryParams: url.Values{
			"request_id": []string{"demo-123"},
		},
	})
	if err != nil {
		log.Fatalf("invoke: %v", err)
	}

	pretty, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(pretty))
}
