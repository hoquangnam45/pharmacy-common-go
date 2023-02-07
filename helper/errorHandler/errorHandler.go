package errorHandler

import "time"

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

func Just[T any](val T) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, error) {
			return val, nil
		},
	}
}

func Error[T any](err error) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, error) {
			var noop T
			return noop, err
		},
	}
}

func Empty[T any]() *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, error) {
			var noop T
			return noop, nil
		},
	}
}

func Factory[T any](f func() T) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, error) {
			val := f()
			return val, nil
		},
	}
}

func FactoryM[T any](f func() (T, error)) *MaybeError[T] {
	return &MaybeError[T]{
		f: func() (T, error) {
			val, err := f()
			return val, err
		},
	}
}

func Lift[T any, K any](f func(T) (K, error)) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return &MaybeError[K]{
			f: func() (K, error) {
				val, err := f(val)
				return val, err
			},
		}
	}
}

func LiftJ[T any, K any](f func(T) K) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return &MaybeError[K]{
			f: func() (K, error) {
				return f(val), nil
			},
		}
	}
}

func LiftN[T any](f func(T) error) func(T) *MaybeError[any] {
	return func(val T) *MaybeError[any] {
		return &MaybeError[any]{
			f: func() (any, error) {
				err := f(val)
				return nil, err
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

func (m *MaybeError[T]) DefaultEval(defaultValue T) T {
	val, err := m.Eval()
	if err != nil {
		return defaultValue
	}
	return val
}

func (m *MaybeError[T]) PanicEval() T {
	val, err := m.Eval()
	if err != nil {
		panic(err)
	}
	return val
}

func (m *MaybeError[T]) EvalWithHandler(handler func(error)) (T, error) {
	val, err := m.Eval()
	if err != nil {
		handler(err)
	}
	return val, err
}

func (m *MaybeError[T]) GoEval(retCh chan<- T, errorCh chan<- error) {
	val, err := m.Eval()
	if err != nil {
		errorCh <- err
	} else {
		retCh <- val
	}
	close(errorCh)
	close(retCh)
}

func (m *MaybeError[T]) EvalWithCh() (<-chan T, <-chan error) {
	out := make(chan T, 1)
	errs := make(chan error, 1)
	go m.GoEval(out, errs)
	return out, errs
}

func (m *MaybeError[T]) Eval() (T, error) {
	if m.hasCache {
		return m.val, m.err
	}
	val, err := m.f()
	if err != nil && m.retry {
		if m.maxRetry > 0 {
			for i := 0; i < m.maxRetry; i++ {
				time.Sleep(m.retryDelay)
				val, err = m.f()
				if err == nil {
					break
				}
			}
		} else if m.maxRetryTime > 0 {
			maxEndTime := time.Now().Add(m.maxRetryTime)
			for time.Now().Before(maxEndTime) {
				time.Sleep(m.retryDelay)
				val, err = m.f()
				if err == nil {
					break
				}
			}
		}
	}
	if m.cache {
		m.val, m.err = val, err
		m.hasCache = true
	}
	return val, err
}

func Peek[T any](f func(T)) func(T) *MaybeError[T] {
	return func(val T) *MaybeError[T] {
		f(val)
		return Just(val)
	}
}

func PeekE[T any](f func(T) error) func(T) *MaybeError[T] {
	return func(val T) *MaybeError[T] {
		err := f(val)
		if err != nil {
			return Error[T](err)
		}
		return Just(val)
	}
}

func FlatMap[T any, K any](m *MaybeError[T], f func(T) *MaybeError[K]) *MaybeError[K] {
	return &MaybeError[K]{
		f: func() (K, error) {
			val, err := m.Eval()
			var noop K
			if err != nil {
				return noop, err
			}
			newVal, newErr := f(val).Eval()
			return newVal, newErr
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
