package lb

import (
	"errors"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/util"
)

var ErrorNeedRefresh = errors.New("reach expire ttl")
var ErrorEmptyList = errors.New("empty list")

type baseLB[T comparable] struct {
	elements    []T
	expireTime  time.Time
	positionMap map[T]int
}

func (l *baseLB[T]) RefreshList(elements map[T]bool, ttl time.Duration) {
	l.elements = util.SetToList(elements)
	l.positionMap = createPositionMap(l.elements)
	l.expireTime = time.Now().Add(ttl)
}

func (l *baseLB[T]) Check() error {
	if len(l.elements) == 0 {
		return ErrorEmptyList
	}
	now := time.Now()
	if now.After(l.expireTime) {
		return ErrorNeedRefresh
	}
	return nil
}

func (l *baseLB[T]) Remove(val T) (int, bool) {
	if idx, ok := l.positionMap[val]; ok {
		l.elements = util.RemoveIndex(l.elements, idx)
		delete(l.positionMap, val)
		return idx, true
	}
	return 0, false
}
