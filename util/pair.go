package util

type Pair[T any, K any] struct {
	First  T
	Second K
}

func NewPair[T any, K any](first T, second K) *Pair[T, K] {
	return &Pair[T, K]{
		First:  first,
		Second: second,
	}
}
