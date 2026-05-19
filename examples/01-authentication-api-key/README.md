# 01 — Authentication with API key

Demonstrates the two ways to wire an API key into a `rest.Client`:

1. **Environment variable.** `NewWithDefaults()` reads `DEEPGRAM_API_KEY` automatically.
2. **Explicit parameter.** `New("your-key", "")` when the key comes from a secret manager / flag / config.

## Run

```bash
export DEEPGRAM_API_KEY=your-key-here
go run ./examples/01-authentication-api-key
```

## Equivalent Python

[`deepgram-python-sdk/examples/01-authentication-api-key.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/01-authentication-api-key.py)
