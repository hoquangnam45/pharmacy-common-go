package errorHandler

func FlatMapV[T any, K any](f func(T) *MaybeError[K], val T) *MaybeError[K] {
	return FlatMap(Just(val), f)
}
