// Example: Live Transcription with WebSocket
//
// Streams a pre-recorded audio file in chunks to simulate a live
// microphone feed. The fixture is 16-bit mono PCM at 44.1 kHz; the
// LiveTranscriptionOptions match that encoding. Real microphone
// integrations use PortAudio bindings or shell out to sox/ffmpeg.
//
// The receive loop type-switches on Event variants. ResultsEvent is
// the only one that carries transcript text; the others signal state
// transitions (speech started, utterance ended, session metadata).

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	wsv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/websocket"
)

const (
	sampleRate = 44100
	chunkSize  = sampleRate * 2
	chunkDelay = 1 * time.Second
)

func main() {
	client := wsv1.NewWithDefaults()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Connect(ctx, &wsv1.LiveTranscriptionOptions{
		Model:      "nova-3",
		Encoding:   "linear16",
		SampleRate: sampleRate,
		Channels:   1,
	})
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer stream.Close()

	fmt.Println("Connection opened")

	_, thisFile, _, _ := runtime.Caller(0)
	audioPath := filepath.Join(filepath.Dir(thisFile), "..", "fixtures", "audio.wav")

	go func() {
		audio, readErr := os.ReadFile(audioPath)
		if readErr != nil {
			log.Printf("read fixture failed: %v", readErr)
			cancel()
			return
		}
		for i := 0; i < len(audio); i += chunkSize {
			end := i + chunkSize
			if end > len(audio) {
				end = len(audio)
			}
			if sendErr := stream.SendAudio(audio[i:end]); sendErr != nil {
				log.Printf("send failed: %v", sendErr)
				return
			}
			time.Sleep(chunkDelay)
		}
		if closeErr := stream.CloseStream(); closeErr != nil {
			log.Printf("CloseStream failed: %v", closeErr)
		}
	}()

	for {
		event, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			fmt.Println("Connection closed")
			return
		}
		if recvErr != nil {
			log.Fatalf("recv failed: %v", recvErr)
		}
		switch e := event.(type) {
		case *wsv1.ResultsEvent:
			if len(e.Channel.Alternatives) > 0 {
				transcript := e.Channel.Alternatives[0].Transcript
				if transcript != "" {
					fmt.Printf("Transcript: %s\n", transcript)
				}
			}
		case *wsv1.MetadataEvent:
			fmt.Println("Received Metadata event")
			return
		case *wsv1.SpeechStartedEvent:
			fmt.Println("Received SpeechStarted event")
		case *wsv1.UtteranceEndEvent:
			fmt.Println("Received UtteranceEnd event")
		case *wsv1.ErrorEvent:
			log.Fatalf("stream error: %s", e.Description)
		}
	}
}
