package errorHandler

import "time"

var emptyFunc func() = func() {}

type MaybeError[T any] struct {
	f            func() (val T, cleanup func(), err error)
	cache        bool
	hasCache     bool
	val          *T
	cleanup      func()
	err          error
	retry        bool
	retryDelay   time.Duration
	maxRetry     int
	maxRetryTime time.Duration
}

func Just[T any](val T) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, func(), error) {
			return val, emptyFunc, nil
		},
	}
}

func Error[T any](err error) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, func(), error) {
			var noop T
			return noop, emptyFunc, err
		},
	}
}

func Empty[T any]() *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, func(), error) {
			var noop T
			return noop, emptyFunc, nil
		},
	}
}

func Transform[T any](f func() (T, error)) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, func(), error) {
			val, err := f()
			return val, emptyFunc, err
		},
	}
}

func Lift[T any, K any](f func(T) (K, error)) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return &MaybeError[K]{
			f: func() (K, func(), error) {
				val, err := f(val)
				return val, emptyFunc, err
			},
		}
	}
}

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

func (m *MaybeError[T]) Cache() *MaybeError[T] {
	return &MaybeError[T]{
		f:        m.f,
		cache:    true,
		hasCache: false,
	}
}

func (m *MaybeError[T]) Cleanup(cleanupT func(T)) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, func(), error) {
			val, cleanup, err := m.Eval()
			if err == nil {
				return val, func() {
					cleanupT(val)
					cleanup()
				}, err
			}
			return val, func() {
				cleanup()
			}, err
		},
	}
}

func (m *MaybeError[T]) Eval() (T, func(), error) {
	if m.hasCache {
		return *m.val, m.cleanup, m.err
	}
	val, cleanup, err := m.f()
	if err != nil {
		if m.retry {
			if m.maxRetry > 0 {
				for i := 0; i < m.maxRetry; i++ {
					time.Sleep(m.retryDelay)
					val, cleanup, err = m.f()
					if err == nil {
						break
					}
					cleanup()
					cleanup = emptyFunc
				}
			} else if m.maxRetryTime > 0 {
				maxEndTime := time.Now().Add(m.maxRetryTime)
				for time.Now().Before(maxEndTime) {
					time.Sleep(m.retryDelay)
					val, cleanup, err = m.f()
					if err == nil {
						break
					}
					cleanup()
					cleanup = emptyFunc
				}
			}
		} else {
			cleanup()
			cleanup = emptyFunc
		}
	}
	if m.cache {
		m.val, m.cleanup, m.err = &val, cleanup, err
		m.hasCache = true
	}
	return val, cleanup, err
}

func FlatMap[T any, K any](m *MaybeError[T], f func(T) *MaybeError[K]) *MaybeError[K] {
	return &MaybeError[K]{
		f: func() (K, func(), error) {
			val, cleanup, err := m.Eval()
			var noop K
			if err != nil {
				return noop, cleanup, err
			}
			newVal, newCleanup, newErr := f(val).Eval()
			cleanup()
			return newVal, newCleanup, newErr
		},
	}
}

func FlatMap2[A any, B any, C any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C]) *MaybeError[C] {
	return FlatMap(FlatMap(m, a), b)
}

func FlatMap3[A any, B any, C any, D any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D]) *MaybeError[D] {
	return FlatMap(FlatMap2(m, a, b), c)
}

func FlatMap4[A any, B any, C any, D any, E any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E]) *MaybeError[E] {
	return FlatMap(FlatMap3(m, a, b, c), d)
}

func FlatMap5[A any, B any, C any, D any, E any, F any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F]) *MaybeError[F] {
	return FlatMap(FlatMap4(m, a, b, c, d), e)
}

func FlatMap6[A any, B any, C any, D any, E any, F any, G any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G]) *MaybeError[G] {
	return FlatMap(FlatMap5(m, a, b, c, d, e), f)
}

func FlatMap7[A any, B any, C any, D any, E any, F any, G any, H any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H]) *MaybeError[H] {
	return FlatMap(FlatMap6(m, a, b, c, d, e, f), g)
}

func FlatMap8[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H],
	h func(H) *MaybeError[I]) *MaybeError[I] {
	return FlatMap(FlatMap7(m, a, b, c, d, e, f, g), h)
}

func FlatMap9[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H],
	h func(H) *MaybeError[I],
	i func(I) *MaybeError[J]) *MaybeError[J] {
	return FlatMap(FlatMap8(m, a, b, c, d, e, f, g, h), i)
}

func FlatMap10[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[J],
	g func(J) *MaybeError[H],
	h func(H) *MaybeError[I],
	i func(I) *MaybeError[J],
	j func(J) *MaybeError[K]) *MaybeError[K] {
	return FlatMap(FlatMap9(m, a, b, c, d, e, f, g, h, i), j)
}
