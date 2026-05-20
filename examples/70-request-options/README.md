# 70 — Request options

`PreRecordedTranscriptionOptions.AdditionalQueryParams` (and the equivalent field on `LiveTranscriptionOptions`) is the escape hatch for passing query parameters the SDK does not yet expose as typed fields. Use it when Deepgram ships a new API parameter you want to test before the SDK is updated.

```go
opts := &restv1.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
    SmartFormat: true,
    AdditionalQueryParams: url.Values{
        "experimental_feature": []string{"true"},
        "custom_tag":           []string{"a", "b"},
    },
}
```

Multiple values per key produce repeated `?key=v1&key=v2` entries. When a key collides with one of the typed fields, `AdditionalQueryParams` wins.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/70-request-options
```

## See also

- [`time` / `context.WithTimeout`](https://pkg.go.dev/context#WithTimeout) — bound a single call.
- [`http.Client`](https://pkg.go.dev/net/http#Client) — pass a custom `*http.Client` via `Client.WithHTTPClient(...)` to add custom headers, configure transports, etc.
