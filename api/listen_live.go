// Code generated from deepgram/spec. DO NOT EDIT.
//
// Source: model/listen-live/listen-live.smithy
// Service: com.deepgram.api.listen.live#ListenLive

package api

// ListenLiveClient is the WSS streaming transcription client.
type ListenLiveClient struct {
	apiKey string
}

// NewListenLiveClient returns a ListenLiveClient authenticated with the given API key.
func NewListenLiveClient(apiKey string) *ListenLiveClient {
	return &ListenLiveClient{apiKey: apiKey}
}

// Stream performs the Stream operation against /v1/listen (WS upgrade).
func (c *ListenLiveClient) Stream() error {
	return nil
}
