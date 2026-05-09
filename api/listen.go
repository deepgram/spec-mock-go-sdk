// Code generated from deepgram/spec. DO NOT EDIT.
//
// Source: model/listen/listen.smithy
// Service: com.deepgram.api.listen#Listen

package api

// ListenClient is the REST batch transcription client.
type ListenClient struct {
	apiKey string
}

// NewListenClient returns a ListenClient authenticated with the given API key.
func NewListenClient(apiKey string) *ListenClient {
	return &ListenClient{apiKey: apiKey}
}

// Transcribe performs the Transcribe operation against /v1/listen.
func (c *ListenClient) Transcribe() error {
	return nil
}
