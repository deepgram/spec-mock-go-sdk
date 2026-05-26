# Error handling

This example shows how to inspect SDK errors with `errors.As`.

HTTP failures are returned as `*httptransport.HTTPError`. When the API
response matches a documented error shape, `HTTPError.Typed` contains a
specific type from `api/types`.

Use this pattern to branch on status codes and typed Deepgram errors
without relying on string matching.

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/71-error-handling
```
