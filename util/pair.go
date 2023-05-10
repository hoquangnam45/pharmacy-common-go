package util

import h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"

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

func PairFM[A any, B any, C any](f func(A) *h.MaybeError[B], anotherValue *h.MaybeError[C]) func(A) *h.MaybeError[*Pair[B, C]] {
	return func(val A) *h.MaybeError[*Pair[B, C]] {
		return h.FlatMap2(
			h.Just(val),
			f,
			func(val B) *h.MaybeError[*Pair[B, C]] {
				anotherVal_, err := anotherValue.Eval()
				if err != nil {
					return h.Error[*Pair[B, C]](err)
				}
				return h.Just(NewPair(val, anotherVal_))
			},
		)
	}
}
