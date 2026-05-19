# 71 — Error handling

Two distinct patterns for two transports, plus context cancellation.

## REST: typed `*HTTPError` via `errors.As`

Every REST response with status >= 400 surfaces as a typed [`*httptransport.HTTPError`](../../api/transport/http/rest.go). Customer code extracts it with `errors.As` and discriminates on `StatusCode`:

```go
var httpErr *httptransport.HTTPError
if errors.As(err, &httpErr) {
    switch httpErr.StatusCode {
    case http.StatusBadRequest:    // 400
    case http.StatusUnauthorized:  // 401
    case http.StatusTooManyRequests: // 429 — read httpErr.Headers.Get("Retry-After")
    }
    log.Printf("body: %s", httpErr.Body)
}
```

Fields available on `*HTTPError`:

| Field | Type | Purpose |
|---|---|---|
| `Method` | `string` | HTTP method of the failed request. |
| `URL` | `string` | Fully-resolved URL including query string. |
| `StatusCode` | `int` | Wire status code. Match against `http.Status*` constants. |
| `Body` | `[]byte` | Raw response body. Decode if you need structured fields. |
| `Headers` | `http.Header` | Response headers. Notably `Retry-After`, `X-Dg-Request-Id`. |
| `Typed` | `error` | Operation-specific decoded error (see below). Nil when the status is not in the operation's declared errors or the body couldn't decode. |

The `.Error()` string format is stable: `"http.Invoke: METHOD URL: STATUS BODY"`. Older code that string-matched the error message keeps working — new code should prefer `errors.As`.

## REST typed errors via `HTTPError.Typed`

For an even sharper handle, `HTTPError.Typed` carries the spec-driven typed error. Each operation declares its possible error types via Smithy `errors: [...]`; the codegen-emitted route metadata includes a decoder that maps the response into one of those types:

```go
var httpErr *httptransport.HTTPError
if errors.As(err, &httpErr) && httpErr.Typed != nil {
    var queryErr *spectypes.InvalidQueryParameterError
    if errors.As(httpErr.Typed, &queryErr) {
        log.Printf("err_code=%s err_msg=%s request_id=%s",
            *queryErr.ErrCode, *queryErr.ErrMsg, *queryErr.RequestId)
        // queryErr.DgError mirrors the dg-error response header
    }

    var rateLimited *spectypes.RateLimitedError
    if errors.As(httpErr.Typed, &rateLimited) {
        // rateLimited.RetryAfter mirrors the Retry-After response header
        time.Sleep(parseRetryAfter(*rateLimited.RetryAfter))
    }
}
```

For Listen REST (`Transcribe` operation), the declared error types and their status codes are:

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

Every type carries the same Legacy-shape body fields — `ErrCode *string`, `ErrMsg *string`, `RequestId *string`, plus `DgError *string` from the `dg-error` response header. `RateLimitedError` also carries `RetryAfter *string` from the `Retry-After` header.

For finer discrimination within a status (e.g. distinguishing `INVALID_QUERY_PARAMETER` from `KEYWORD_LIMIT_EXCEEDED` within 400), inspect the `ErrCode` field of the typed struct. The full set of known err_code values is documented in [`deepgram/spec/model/common/primitives.smithy`](https://github.com/deepgram/spec/blob/main/model/common/primitives.smithy) on the `ErrCode` type.

When the status isn't in the operation's declared list (e.g. an upstream `502 Bad Gateway` returns an HTML body), `Typed` is `nil` — `Body` and `StatusCode` are still populated so customer code can log the raw failure.

## WebSocket: `*wsv1.ErrorEvent` via Recv-loop type-switch

Server-emitted WebSocket errors arrive as an `Event` variant on `stream.Recv()`, not as exceptions:

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
        // handle transcript ...
    }
}
```

The `ErrorEvent.Description` carries Deepgram's text; the `Type` field carries the discriminator from the wire (e.g. `"Error"`).

## WebSocket: facade-level sentinels via `errors.Is`

Send-side validation errors are sentinels — they fire BEFORE bytes hit the wire, so the customer can catch them locally:

| Sentinel | When |
|---|---|
| [`wsv1.ErrFrameTooLarge`](../../pkg/client/listen/v1/websocket/resilience.go) | `SendAudio` chunk exceeds `Config.MaxFrameSizeBytes`. |
| [`wsv1.ErrSendTimeout`](../../pkg/client/listen/v1/websocket/resilience.go) | `SendAudio` blocked longer than `Config.SendTimeout`. |

Match with `errors.Is(err, wsv1.ErrFrameTooLarge)`.

## Context cancellation

Any call accepting a `context.Context` honours its cancellation. `context.WithTimeout` is the canonical way to bound a single call:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := client.FromURL(ctx, url, opts)
if errors.Is(err, context.DeadlineExceeded) { ... }
```

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/71-error-handling
```

The example runs four sub-cases. Sections 1 + 2 exercise a real 400 response; section 3 forces a 1ms deadline so context cancellation fires; section 4 prints the WebSocket-error pattern as a code template (without standing up a real WS connection).

## Equivalent Python

[`deepgram-python-sdk/examples/71-error-handling.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/71-error-handling.py) teaches the Python SDK's exception hierarchy (`ApiError` → `BadRequestError`) plus `EventType.ERROR` callbacks for WebSocket.

| Python | Go |
|---|---|
| `except ApiError as e: e.status_code, e.body, e.headers` | `var httpErr *httptransport.HTTPError; errors.As(err, &httpErr); use httpErr.StatusCode, httpErr.Body, httpErr.Headers` |
| `except BadRequestError as e` | `if httpErr.StatusCode == http.StatusBadRequest` |
| `except Exception as e` (network catch-all) | `errors.As` returns false → handle as transport error |
| `connection.on(EventType.ERROR, on_error)` callback | `case *wsv1.ErrorEvent:` arm of the Recv-loop type-switch |

### Closed gap (was a known gap in v0.0.x)

Earlier versions of this SDK shipped only the generic `*HTTPError` container with no per-status typed sub-errors — customers had to discriminate on `httpErr.StatusCode` and decode `httpErr.Body` themselves. That gap is now closed: each operation's Smithy spec declares `errors: [...]`, the codegen emits a per-operation decoder, and `HTTPError.Typed` carries the typed error directly. See the "REST typed errors via `HTTPError.Typed`" section above.
