package lb

import (
	"time"
)

type RoundRobinLB[T comparable] struct {
	*baseLB[T]
	idx int
	LoadBalancer[T]
}

func NewRoundRobinLB[T comparable](elementFetcher ElementFetcher[T], ttl time.Duration) *RoundRobinLB[T] {
	return &RoundRobinLB[T]{
		baseLB: NewBaseLb(elementFetcher, ttl),
	}
}

func (l *RoundRobinLB[T]) LoadBalancing() T {
	ret := l.elements[l.idx]
	l.idx = (l.idx + 1) % len(l.elements)
	return ret
}

func (l *RoundRobinLB[T]) Get() (T, error) {
	var noop T
	if err := l.CheckRefresh(); err != nil {
		return noop, err
	}
	return l.LoadBalancing(), nil
}
