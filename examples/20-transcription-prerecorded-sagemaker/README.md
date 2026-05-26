# Prerecorded transcription with SageMaker

This example routes prerecorded transcription through SageMaker by
passing `prerecorded.WithSageMakerTransport` to `prerecorded.New`.

The standard typed options still become the request payload. When
`AdditionalQueryParams` is set, this transport encodes those values into
the SageMaker `CustomAttributes` header.

Use `WithTargetVariant`, `WithTargetModel`, `WithInferenceID`, and
`WithEnableExplanations` when your endpoint needs those settings.
