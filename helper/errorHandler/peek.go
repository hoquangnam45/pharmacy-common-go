package errorHandler

func PeekE[T any](f func(T) error) func(T) *MaybeError[T] {
	return func(val T) *MaybeError[T] {
		err := f(val)
		if err != nil {
			return Error[T](err)
		}
		return Just(val)
	}
}

func PeekEVa[T any](f func(...T) error) func(T) *MaybeError[T] {
	return PeekE(func(val T) error {
		return f(val)
	})
}
