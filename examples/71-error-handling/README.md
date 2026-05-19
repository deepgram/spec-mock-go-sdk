# 71 — Error handling

Patterns for handling errors from the SDK.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/71-error-handling
```

## REST

```go
var httpErr *httptransport.HTTPError
if errors.As(err, &httpErr) {
    var queryErr *spectypes.InvalidQueryParameterError
    if errors.As(httpErr.Typed, &queryErr) {
        // *queryErr.ErrCode, *queryErr.ErrMsg, *queryErr.RequestId, *queryErr.DgError
    }

    var rateLimited *spectypes.RateLimitedError
    if errors.As(httpErr.Typed, &rateLimited) {
        // *rateLimited.RetryAfter from the Retry-After response header
    }

    if httpErr.Typed == nil {
        // Status not in the operation's declared errors, or body not JSON.
        // httpErr.Body / httpErr.StatusCode / httpErr.Headers still populated.
    }
}
```

Declared error types for `Transcribe` (`POST /v1/listen`):

| Type | Status |
|---|---|
| `*spectypes.InvalidQueryParameterError` | 400 |
| `*spectypes.UnauthorizedError` | 401 |
| `*spectypes.PaymentRequiredError` | 402 |
| `*spectypes.ForbiddenError` | 403 |
| `*spectypes.NotFoundError` | 404 |
| `*spectypes.SlowUploadError` | 408 |
| `*spectypes.PayloadTooLargeError` | 413 |
| `*spectypes.UnsupportedMediaTypeError` | 415 |
| `*spectypes.RateLimitedError` | 429 |
| `*spectypes.InternalServerError` | 500 |

Fine-grained discrimination within a status: read `*typed.ErrCode`.

## Context cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := client.FromURL(ctx, url, opts)
if errors.Is(err, context.DeadlineExceeded) { ... }
```

## WebSocket

Server errors via `stream.Recv()`:

```go
for {
    event, err := stream.Recv()
    if errors.Is(err, io.EOF) { return }
    if err != nil { return err }
    switch e := event.(type) {
    case *wsv1.ErrorEvent:
        log.Printf("deepgram error: %s", e.Description)
        return
    case *wsv1.ResultsEvent:
        // ...
    }
}
```

Client-level send-side sentinels:

| Sentinel | When |
|---|---|
| `wsv1.ErrFrameTooLarge` | `SendAudio` chunk exceeds `Config.MaxFrameSizeBytes` |
| `wsv1.ErrSendTimeout` | `SendAudio` blocked longer than `Config.SendTimeout` |
