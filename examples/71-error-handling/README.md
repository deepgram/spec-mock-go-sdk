# 71 — Error handling

Go errors are plain values. The SDK uses the standard idioms:

- `errors.Is(err, sentinel)` — match against an exported sentinel.
- `errors.As(err, &target)` — type-assert through wrapping.
- `context` cancellation — pass a `context.WithTimeout` / `context.WithCancel` to bound the call.

The SDK does **not** have a custom Error class hierarchy; everything comes back as a `error` with the relevant status / body / wrapped cause embedded in the message.

## Exported sentinel errors

| Sentinel | Where | Match with |
|---|---|---|
| [`wsv1.ErrFrameTooLarge`](../../pkg/client/listen/v1/websocket/resilience.go) | WS `SendAudio` chunk exceeds `Config.MaxFrameSizeBytes` | `errors.Is(err, wsv1.ErrFrameTooLarge)` |
| [`wsv1.ErrSendTimeout`](../../pkg/client/listen/v1/websocket/resilience.go) | WS `SendAudio` blocked longer than `Config.SendTimeout` | `errors.Is(err, wsv1.ErrSendTimeout)` |
| `context.DeadlineExceeded` | Any call where the supplied `ctx` deadline fires | `errors.Is(err, context.DeadlineExceeded)` |
| `context.Canceled` | Any call where the supplied `ctx` is `cancel()`'d | `errors.Is(err, context.Canceled)` |

## What this example covers

1. **Bad URL** — REST surfaces the underlying DNS / transport error.
2. **Empty URL** — Deepgram returns 400; the SDK wraps the response body in the error.
3. **Context cancellation** — `context.WithTimeout` then `errors.Is(err, context.DeadlineExceeded)`.
4. **WebSocket sentinel** — `wsv1.ErrFrameTooLarge` from `SendAudio` when a chunk exceeds the configured limit.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/71-error-handling
```

## Equivalent Python

[`deepgram-python-sdk/examples/71-error-handling.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/71-error-handling.py)

| Python | Go |
|---|---|
| `except ApiError as e: e.status_code, e.body` | Inspect the error string, or wrap calls in your own `*HTTPError` type and unwrap. |
| `except BadRequestError as e` | No specific BadRequest type; the SDK returns the wrapped status. |
| `try/finally` for cleanup | `defer stream.Close()`. |

## Best practices

1. Always pass a `context.Context` you control. Bound long-running calls with `context.WithTimeout`.
2. Match exported sentinels with `errors.Is`. Do not string-match error messages.
3. Defer cleanup (`stream.Close`, `resp.Body.Close`) so panics still tidy up.
4. For retries: the SDK does not auto-retry. Wrap the call site in your own retry loop and decide what's retryable based on the error.
