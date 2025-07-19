package kafka_retry

import (
	"errors"
	"strings"
)

type RetryableError struct {
	Err error
	Msg string
}

func (e *RetryableError) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func NewRetryableError(err error, msg string) *RetryableError {
	return &RetryableError{Err: err, Msg: msg}
}

type NonRetryableError struct {
	Err error
	Msg string
}

func (e *NonRetryableError) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func NewNonRetryableError(err error, msg string) *NonRetryableError {
	return &NonRetryableError{Err: err, Msg: msg}
}

func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	var nonRetryable *NonRetryableError
	if errors.As(err, &nonRetryable) {
		return false
	}

	var retryable *RetryableError
	if errors.As(err, &retryable) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	retryableKeywords := []string{
		"connection",
		"timeout",
		"network",
		"temporary",
		"unavailable",
		"deadlock",
	}

	for _, keyword := range retryableKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}
