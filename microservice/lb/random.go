package lb

import (
	"math/rand"
)

type RandomLB[T comparable] struct {
	*baseLB[T]
}

func NewRandomLB[T comparable]() *RandomLB[T] {
	return &RandomLB[T]{
		baseLB: &baseLB[T]{},
	}
}

func (l *RandomLB[T]) Get() (T, error) {
	err := l.Check()
	var noop T
	if err != nil {
		return noop, err
	}
	ret := l.elements[rand.Intn(len(l.elements))]
	return ret, nil
}
