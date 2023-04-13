package errorHandler

import "time"

func (m *MaybeError[T]) RetryUntilSuccess(maxRetryTime time.Duration, retryDelay time.Duration) *MaybeError[T] {
	return &MaybeError[T]{
		f:            m.f,
		retry:        true,
		maxRetryTime: maxRetryTime,
		retryDelay:   retryDelay,
	}
}

func (m *MaybeError[T]) Retry(maxRetry int, retryDelay time.Duration) *MaybeError[T] {
	return &MaybeError[T]{
		f:          m.f,
		retry:      true,
		maxRetry:   maxRetry,
		retryDelay: retryDelay,
	}
}

func RetryFM[T any, K any](f func(T) *MaybeError[K], maxRetry int, retryDelay time.Duration) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return FlatMap(Just(val), f).Retry(maxRetry, retryDelay)
	}
}

func RetryUntilSuccessFM[T any, K any](f func(T) *MaybeError[K], maxRetryTime time.Duration, retryDelay time.Duration) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return FlatMap(Just(val), f).RetryUntilSuccess(maxRetryTime, retryDelay)
	}
}
