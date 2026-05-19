# 02 — Authentication with access token

Access tokens are short-lived bearer credentials, typically issued by an upstream identity provider (in Deepgram's case `dx-id`). This SDK does not include a token-issuance helper — issuance is a separate concern owned by whatever auth flow the customer's application uses.

When both `DEEPGRAM_API_KEY` and `DEEPGRAM_ACCESS_TOKEN` are set, the SDK prefers the access token.

## Run

```bash
export DEEPGRAM_ACCESS_TOKEN=your-token-here
go run ./examples/02-authentication-access-token
```

## Equivalent Python

[`deepgram-python-sdk/examples/02-authentication-access-token.py`](https://github.com/deepgram/deepgram-python-sdk/blob/main/examples/02-authentication-access-token.py)

The Python SDK ships an `auth.v1.tokens.grant()` helper that calls Deepgram's token-grant endpoint. This SDK currently does not — token issuance is upstream.
