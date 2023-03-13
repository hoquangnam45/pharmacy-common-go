package errorHandler

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

func LiftVa[T any, K any](f func(...T) (K, error)) func(T) *MaybeError[K] {
	return Lift(func(val T) (K, error) {
		return f(val)
	})
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

func LiftJVa[T any, K any](f func(...T) K) func(T) *MaybeError[K] {
	return LiftJ(func(val T) K {
		return f(val)
	})
}

func LiftE[T any](f func(T) error) func(T) *MaybeError[any] {
	return func(val T) *MaybeError[any] {
		return &MaybeError[any]{
			f: func() (any, error) {
				err := f(val)
				return nil, err
			},
		}
	}
}

func LiftEVa[T any](f func(...T) error) func(T) *MaybeError[any] {
	return LiftE(func(val T) error {
		return f(val)
	})
}

func LiftFactoryE[T any](f func() error) func(T) *MaybeError[any] {
	return LiftE(func(val T) error {
		return f()
	})
}
