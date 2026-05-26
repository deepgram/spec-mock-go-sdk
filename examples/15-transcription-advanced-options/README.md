# Advanced transcription options

This example combines commonly used typed options for prerecorded
transcription.

It also includes `AdditionalQueryParams`, which lets you send a query
parameter before the SDK exposes a typed field for it.

Use this pattern when most options are already typed and you need one
temporary escape hatch.

```bash
export DEEPGRAM_API_KEY="your_api_key"
go run ./examples/15-transcription-advanced-options
```
