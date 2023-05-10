package lb

import (
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

type RoundRobinLB[T comparable] struct {
	*baseLB[T]
	idx int
}

func NewRoundRobinLB[T comparable](elementFetcher ElementFetcher[T], ttl time.Duration, logger log.Logger) LoadBalancer[T] {
	return &RoundRobinLB[T]{
		baseLB: NewBaseLb(elementFetcher, ttl, logger),
	}
}

func (l *RoundRobinLB[T]) LoadBalancing() (T, error) {
	ret := l.elements[l.idx]
	l.idx = (l.idx + 1) % len(l.elements)
	return ret, nil
}

func (l *RoundRobinLB[T]) Get() (T, error) {
	var noop T
	if err := l.CheckRefresh(); err != nil {
		return noop, err
	}
	return l.LoadBalancing()
}
