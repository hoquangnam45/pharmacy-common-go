package errorHandler

import (
	"time"
)

type MaybeError[T any] struct {
	f            func() (val T, err error)
	cache        bool
	hasCache     bool
	val          T
	err          error
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
