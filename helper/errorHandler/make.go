package errorHandler

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

func FactoryE(f func() error) *MaybeError[any] {
	return &MaybeError[any]{
		f: func() (any, error) {
			err := f()
			return nil, err
		},
	}
}
