package lb

import (
	"math/rand"
	"time"
)


type RandomLB[T comparable] struct {
	*baseLB[T]
	LoadBalancer[T]
}

func NewRandomLB[T comparable](elementFetcher ElementFetcher[T], ttl time.Duration) *RandomLB[T] {
	return &RandomLB[T]{
		baseLB: NewBaseLb(elementFetcher, ttl),
	}
}

func (l *RandomLB[T]) LoadBalancing() (T, error) {
	var noop T
	if err := l.Check(); err != nil {
		return noop, err
	}
	ret := l.elements[rand.Intn(len(l.elements))]
	return ret, nil
}

func (l *RandomLB[T]) Get() (T, error) {
	var noop T
	if err := l.CheckRefresh(); err != nil {
		return noop, err
	}
	return l.LoadBalancing()
}
