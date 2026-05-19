// Example: Authentication with Access Token
//
// Access-token auth via DEEPGRAM_ACCESS_TOKEN env var or explicit
// New("", "<token>"). Access tokens are short-lived bearer credentials,
// typically issued by an upstream identity provider (in Deepgram's case
// dx-id). This SDK does not include a token-issuance helper — the
// assumption is the customer's application already has a token from
// wherever it manages auth.
//
// When both API key and access token are configured, the SDK prefers
// the access token.

package main

import (
	"fmt"
	"os"

	restv1 "github.com/deepgram/spec-mock-go-sdk/pkg/client/listen/v1/rest"
)

func main() {
	if os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		fmt.Println("Set DEEPGRAM_ACCESS_TOKEN to run this example end-to-end.")
		fmt.Println("Access tokens are typically issued by dx-id; see")
		fmt.Println("https://id.dx.deepgram.com for the provisioning surface.")
		return
	}

	_ = restv1.NewWithDefaults()
	fmt.Println("Client initialized with access token from DEEPGRAM_ACCESS_TOKEN")

	_ = restv1.New("", "your-access-token-here")
	fmt.Println("Also constructible via restv1.New(\"\", \"<token>\")")
}
