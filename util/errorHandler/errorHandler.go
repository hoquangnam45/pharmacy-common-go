package errorHandler

import (
	"time"
)

type result[T any] struct {
	val T
	err error
}

type MaybeError[T any] struct {
	f            func() (val T, err error)
	cache        bool
	hasCache     bool
	result       *result[T]
	retry        bool
	retryDelay   time.Duration
	maxRetry     int
	maxRetryTime time.Duration
}

func (m *MaybeError[T]) Unwrap() func() (T, error) {
	return func() (T, error) {
		return m.Eval()
	}
}
