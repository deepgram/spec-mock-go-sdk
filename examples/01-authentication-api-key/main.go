// Example: Authentication with API Key
//
// API-key auth via DEEPGRAM_API_KEY env var or explicit New(...). The
// SDK reads DEEPGRAM_API_KEY when you call NewWithDefaults; pass the
// key directly to New when you want to wire it from somewhere else
// (secret manager, CLI flag, config file).

package main

import (
	"fmt"
	"os"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	client := restv1.NewWithDefaults()
	fmt.Println("Client initialized with API key:", client != nil)

	if explicit := os.Getenv("DEEPGRAM_API_KEY"); explicit != "" {
		_ = restv1.New(explicit, "")
		fmt.Println("Also constructible via restv1.New(\"<key>\", \"\")")
	}
}
