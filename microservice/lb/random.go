package lb

import (
	"math/rand"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

type RandomLB[T comparable] struct {
	*baseLB[T]
}

func NewRandomLB[T comparable](elementFetcher ElementFetcher[T], ttl time.Duration, logger log.Logger) LoadBalancer[T] {
	return &RandomLB[T]{
		baseLB: NewBaseLb(elementFetcher, ttl, logger),
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
