package lb

import (
	"errors"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/util"
)

var ErrorNeedRefresh = errors.New("reach expire ttl")
var ErrorEmptyList = errors.New("empty list")

type ElementFetcher[T comparable] func() (map[T]bool, error)

type LoadBalancer[T comparable] interface {
	LoadBalancing() (T, error)
	Get() (T, error)
}

type baseLB[T comparable] struct {
	elements       []T
	ttl            time.Duration
	activeTime     time.Time
	positionMap    map[T]int
	elementFetcher ElementFetcher[T]
}

func NewBaseLb[T comparable](elementFetcher ElementFetcher[T], ttl time.Duration) *baseLB[T] {
	return &baseLB[T]{
		elementFetcher: elementFetcher,
		ttl:            ttl,
	}
}

func (l *baseLB[T]) RefreshList() error {
	newElements, err := l.elementFetcher()
	if err != nil {
		return err
	}
	l.elements = util.SetToList(newElements)
	l.positionMap = createPositionMap(l.elements)
	l.activeTime = time.Now()
	return nil
}

func (l *baseLB[T]) Check() error {
	if len(l.elements) == 0 {
		return ErrorEmptyList
	}
	if time.Since(l.activeTime) >= l.ttl {
		return ErrorNeedRefresh
	}
	return nil
}

func (l *baseLB[T]) CheckRefresh() error {
	err := l.Check()
	if err != nil {
		err := l.RefreshList()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *baseLB[T]) Remove(val T) (int, bool) {
	if idx, ok := l.positionMap[val]; ok {
		util.SugaredLogger.Infof("Delete from LB at index %d  for lengths %d", idx, len(l.elements))
		l.elements = util.RemoveIndex(l.elements, idx)
		delete(l.positionMap, val)
		return idx, true
	}
	return 0, false
}

func (l *baseLB[T]) Add(val T) {
	if _, ok := l.positionMap[val]; ok {
		return
	}
	if l.elements == nil {
		l.elements = []T{}
	}
	l.elements = append(l.elements, val)
	l.positionMap[val] = len(l.elements) - 1
}

func (l *baseLB[T]) Exist(val T) bool {
	if _, ok := l.positionMap[val]; ok {
		return true
	}
	return false
}
