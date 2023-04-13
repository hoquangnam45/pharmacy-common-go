package errorHandler

import "github.com/hoquangnam45/pharmacy-common-go/util"

func PairFM[A any, B any, C any](f func(A) *MaybeError[B], anotherValue *MaybeError[C]) func(A) *MaybeError[*util.Pair[B, C]] {
	return func(val A) *MaybeError[*util.Pair[B, C]] {
		return FlatMap2(
			Just(val),
			f,
			func(val B) *MaybeError[*util.Pair[B, C]] {
				anotherVal_, err := anotherValue.Eval()
				if err != nil {
					return Error[*util.Pair[B, C]](err)
				}
				return Just(util.NewPair(val, anotherVal_))
			},
		)
	}
}
