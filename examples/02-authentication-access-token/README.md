# 02 — Authentication with access token

Access tokens are short-lived bearer credentials. This SDK does not include a token-issuance helper — your application is expected to obtain a token through whatever auth flow it already uses.

When both `DEEPGRAM_API_KEY` and `DEEPGRAM_ACCESS_TOKEN` are set, the SDK prefers the access token.

## Run

```bash
export DEEPGRAM_ACCESS_TOKEN=your-token-here
go run ./examples/02-authentication-access-token
```
