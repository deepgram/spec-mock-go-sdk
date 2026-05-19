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

The `.Error()` string format is stable: `"http.Invoke: METHOD URL: STATUS BODY"`. Older code that string-matched the error message keeps working — new code should prefer `errors.As`.

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

### Known gap (intentional)

The Python SDK has Fern-generated typed subclasses for each documented status code (`BadRequestError`, future `UnauthorizedError`, etc.). This Go SDK currently has the generic `*HTTPError` container; spec-driven per-status typed errors are a follow-up that requires adding `@error` + `@httpError(N)` structures to the Smithy spec and extending the codegen to decode response bodies into them.

Until then, customers discriminate on `httpErr.StatusCode` rather than on Go type. Status-code switching is structurally identical and works against the typed error you already have.
