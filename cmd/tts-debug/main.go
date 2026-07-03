// tts-debug hits the TTS SageMaker endpoint via the raw transport (no resilient
// retry masking) to surface the actual error from a Speak message.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	sm "github.com/deepgram/spec-mock-go-sdk/api/transport/sagemaker"
	spectypes "github.com/deepgram/spec-mock-go-sdk/api/types"
)

func main() {
	endpoint := flag.String("endpoint", "dg-tts-test", "")
	region := flag.String("region", "us-east-2", "")
	model := flag.String("model", "aura-2-atlas-en", "")
	query := flag.String("query", "", "extra query (e.g. encoding=linear16&sample_rate=16000)")
	flag.Parse()

	ctx := context.Background()
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	client := sm.NewBidiClient(cfg, sm.DefaultConfig())
	q := "model=" + *model
	if *query != "" {
		q += "&" + *query
	}
	fmt.Printf("path=v1/speak query=%q\n", q)
	stream, err := sm.OpenStream[spectypes.SpeakClientStream, spectypes.SpeakServerStream](
		ctx, client, *endpoint, "v1/speak", q, "", "", "", "",
		spectypes.MarshalSpeakClientStream, spectypes.UnmarshalSpeakServerStream)
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		os.Exit(1)
	}
	defer stream.Close()

	go func() {
		for {
			m, err := stream.Recv()
			if err != nil {
				fmt.Println("RECV ERR:", err)
				return
			}
			fmt.Printf("RECV: %T\n", m)
		}
	}()

	t := "Hello world, this is a self hosted text to speech smoke test."
	if err := stream.Send(&spectypes.SpeakClientStreamMemberSpeak{Value: spectypes.SpeakText{Text: &t}}); err != nil {
		fmt.Println("SEND ERR:", err)
	}
	if err := stream.Send(&spectypes.SpeakClientStreamMemberFlush{Value: spectypes.SpeakFlush{}}); err != nil {
		fmt.Println("FLUSH ERR:", err)
	}
	time.Sleep(8 * time.Second)
}
