# Additional request options

This example demonstrates `AdditionalQueryParams`.

Use it for parameters the API supports before the SDK has a typed field.
If a key collides with a typed field, the additional value wins.

The example sets `Model` and then overrides the `model` query parameter
through `AdditionalQueryParams` to show the precedence rule.

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/70-request-options
```
