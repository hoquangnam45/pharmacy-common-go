package errorHandler

func (m *MaybeError[T]) Cache() *MaybeError[T] {
	return &MaybeError[T]{
		f:        m.f,
		cache:    true,
		hasCache: false,
	}
}

func CacheFM[T any, K any](f func(T) *MaybeError[K]) func(T) *MaybeError[K] {
	return func(val T) *MaybeError[K] {
		return FlatMap(Just(val), f).Cache()
	}
}
