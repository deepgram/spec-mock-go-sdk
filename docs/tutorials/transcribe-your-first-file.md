# Tutorial: Transcribe your first file

End-to-end walk-through of the listen REST surface. Sends one audio file to `POST /v1/listen` and prints the transcript.

By the end you will:

- Have a runnable Go program that talks to Deepgram.
- Understand the per-package layout (`options.go` + `convert.go` + `client.go`).
- Know which options matter for typical pre-recorded transcription.

## Prerequisites

- Go 1.24 or later.
- A Deepgram API key. Get one at [console.deepgram.com](https://console.deepgram.com).
- An audio file URL accessible over HTTPS. For this tutorial we use [the Bueller sample](https://dpgr.am/bueller.wav) Deepgram hosts publicly.

## 1. Set up the project

```bash
mkdir transcribe-demo && cd transcribe-demo
go mod init example.com/transcribe-demo
go get github.com/deepgram/spec-mock-go-sdk@latest
```

Once the module path renames per [ADR-0001](../adrs/0001-module-path.md), substitute `github.com/deepgram/sdk-go@latest`.

Set your API key:

```bash
export DEEPGRAM_API_KEY=your-key-here
```

## 2. Write the program

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()

	result, err := client.FromURL(
		context.Background(),
		"https://dpgr.am/bueller.wav",
		&restv1.PreRecordedTranscriptionOptions{
			Model:       "nova-3",
			SmartFormat: true,
			Punctuate:   true,
			Language:    "en-US",
		},
	)
	if err != nil {
		log.Fatalf("transcribe failed: %v", err)
	}

	pretty, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(pretty))
	os.Exit(0)
}
```

Save as `main.go`.

## 3. Run it

```bash
go run main.go
```

You should see a JSON document with `Results.Channels[0].Alternatives[0].Transcript` containing the spoken text from the audio file.

## What the code does

| Line | What it does |
|---|---|
| `restv1.NewWithDefaults()` | Reads `DEEPGRAM_API_KEY` (or `DEEPGRAM_ACCESS_TOKEN`) from the environment. |
| `&PreRecordedTranscriptionOptions{...}` | Builds the customer-facing options. Every field corresponds to one `@httpQuery` parameter the Deepgram REST API accepts. Zero values are omitted from the wire. |
| `client.FromURL(ctx, url, opts)` | Sends a `POST /v1/listen` with a `{"url": "..."}` JSON body and the options serialised as query string. Returns the parsed response. |
| `Result` | The full response shape, defined in `api/types.go` and re-exported from the facade. |

## Common option choices

| Option | When to set it |
|---|---|
| `Model` | Always. `"nova-3"` is the current general-purpose model. Pick a specialised model (`"nova-2-meeting"`, `"nova-2-phonecall"`, etc.) for specific audio types. |
| `Language` | When the audio is not English, or when you want to pin the language explicitly. Format: BCP-47 (`"es"`, `"fr-CA"`, etc.). |
| `SmartFormat` | Almost always `true`. Formats dates, numbers, addresses naturally instead of as raw speech. |
| `Punctuate` | Almost always `true`. Inserts periods, commas, question marks. |
| `Diarize` | When you need speaker-separated output. Adds a `speaker` field to every word. |
| `Utterances` | When you want utterance-level grouping (transcript broken into spoken phrases). |
| `Paragraphs` | When you want paragraph-level grouping. Higher-level than utterances. |
| `Redact` | When you need PII / PCI / SSN redaction. Pass `[]string{"pci", "ssn", "numbers"}` etc. |

The full inventory lives in [`pkg/client/listen/v1/rest/options.go`](../../pkg/client/listen/v1/rest/options.go).

## Sending audio bytes (not a URL)

For audio you have in-memory:

```go
audio, _ := os.ReadFile("clip.wav")
result, err := client.FromStream(
	context.Background(),
	bytes.NewReader(audio),
	"audio/wav",
	&restv1.PreRecordedTranscriptionOptions{Model: "nova-3", SmartFormat: true},
)
```

`FromFile` is a shorthand for the local-path case.

## Handling errors

`FromURL` / `FromFile` / `FromStream` return the parsed `Result` and a Go `error`. The error is non-nil for:

- HTTP transport failures (DNS, TLS, network).
- Non-2xx HTTP responses from Deepgram (auth failure, malformed request, rate-limit).
- Response-decoding failures (the SDK couldn't parse the response body).

The error message includes the HTTP status and Deepgram's error body when present. Wrap the call with `errors.Is` / `errors.As` if you need to branch on specific failures.

## Where to go from here

- Add `Callback` to make the request async (Deepgram POSTs back to your URL when done).
- Migrate to live streaming via the [WebSocket tutorial](./build-a-voice-agent.md) when you need real-time.
- Read [`pkg/client/listen/v1/rest/example_test.go`](../../pkg/client/listen/v1/rest/example_test.go) — every public symbol on the REST surface has a runnable godoc example.
