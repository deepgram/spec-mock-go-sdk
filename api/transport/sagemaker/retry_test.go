// Copyright Deepgram, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

// Unit tests for the pure resilience algorithms ported from the burst-tested
// Java transport (deepgram-sagemaker-java): error classification and
// exponential-backoff-with-full-jitter. These need no AWS endpoint.
package sagemaker

import (
	"errors"
	"testing"

	"github.com/aws/smithy-go"
)

// apiErr is a minimal smithy.APIError + HTTPStatusCode carrier, matching how
// aws-sdk-go-v2 surfaces service errors.
type apiErr struct {
	code   string
	status int
}

func (e apiErr) Error() string                 { return e.code }
func (e apiErr) ErrorCode() string             { return e.code }
func (e apiErr) ErrorMessage() string          { return e.code }
func (e apiErr) ErrorFault() smithy.ErrorFault { return smithy.FaultUnknown }
func (e apiErr) HTTPStatusCode() int           { return e.status }

func TestClassify(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want Classification
	}{
		{"throttling-coded 400 is retryable", apiErr{"ThrottlingException", 400}, Retryable},
		{"403 forbidden is terminal", apiErr{"AccessDeniedException", 403}, Terminal},
		{"404 not found is terminal", apiErr{"ResourceNotFound", 404}, Terminal},
		{"400 validation is terminal", apiErr{"ValidationException", 400}, Terminal},
		{"429 rate-limit is retryable", apiErr{"TooManyRequests", 429}, Retryable},
		{"424 failed-dependency is retryable", apiErr{"ModelError", 424}, Retryable},
		{"500 server error is retryable", apiErr{"InternalFailure", 500}, Retryable},
		{"non-API transient error is retryable", errors.New("connection reset by peer"), Retryable},
		{"nil is retryable", nil, Retryable},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.err); got != c.want {
				t.Fatalf("Classify(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

func TestClassifyUnwrapsWrappedError(t *testing.T) {
	wrapped := errors.Join(errors.New("outer"), apiErr{"ValidationException", 400})
	if got := Classify(wrapped); got != Terminal {
		t.Fatalf("Classify(wrapped 400) = %v, want Terminal", got)
	}
}

func TestComputeBackoff(t *testing.T) {
	const initial, max, mult = int64(100), int64(5000), 2.0

	// Jitter floored to the low end (randIntn returns 0) yields exactly initialMs.
	zero := func(n int64) int64 { return 0 }
	if got := ComputeBackoff(initial, max, mult, 0, zero); got != initial {
		t.Fatalf("attempt 0 (jitter=0) = %d, want %d", got, initial)
	}

	// Jitter at the high end (randIntn returns n-1) yields the ceiling. At
	// attempt 2 the exponential ceiling is 100*2^2 = 400.
	high := func(n int64) int64 { return n - 1 }
	if got := ComputeBackoff(initial, max, mult, 2, high); got != 400 {
		t.Fatalf("attempt 2 (jitter=max) = %d, want 400", got)
	}

	// The ceiling is capped at maxMs: 100*2^10 = 102400 -> capped to 5000.
	if got := ComputeBackoff(initial, max, mult, 10, high); got != max {
		t.Fatalf("attempt 10 (capped) = %d, want %d", got, max)
	}

	// Every result stays within [initialMs, maxMs] regardless of attempt.
	mid := func(n int64) int64 { return n / 2 }
	for attempt := 0; attempt < 20; attempt++ {
		got := ComputeBackoff(initial, max, mult, attempt, mid)
		if got < initial || got > max {
			t.Fatalf("attempt %d = %d, out of [%d,%d]", attempt, got, initial, max)
		}
	}
}

func TestDefaultConfigValues(t *testing.T) {
	c := DefaultConfig()
	if c.MaxStreamsPerConnection != 1 {
		t.Fatalf("MaxStreamsPerConnection = %d, want 1 (one stream per connection)", c.MaxStreamsPerConnection)
	}
	if c.MaxConcurrency != 500 || c.MaxRetries != 5 {
		t.Fatalf("unexpected concurrency/retries: %+v", c)
	}
	if c.ConnectionTimeout.Seconds() != 30 || c.ConnectionAcquireTimeout.Seconds() != 60 {
		t.Fatalf("burst-tuned timeouts not applied: %+v", c)
	}
}
