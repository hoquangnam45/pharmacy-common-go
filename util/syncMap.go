package util

import "sync"

type SyncMap[T comparable, K any] struct {
	sync.RWMutex
	data map[T]K
}

func WrapToSyncMap[T comparable, K any](m map[T]K) *SyncMap[T, K] {
	return &SyncMap[T, K]{
		data: m,
	}
}

func NewSyncMap[T comparable, K any]() *SyncMap[T, K] {
	return &SyncMap[T, K]{
		data: map[T]K{},
	}
}

func (s *SyncMap[T, K]) Get(key T) K {
	s.RLock()
	defer s.RUnlock()
	return s.data[key]
}

func (s *SyncMap[T, K]) GetOrDefault(key T, defaultValue K) K {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	if ok {
		return val
	}
	return defaultValue
}

func (s *SyncMap[T, K]) GetOrSet(key T, defaultValue K) K {
	s.Lock()
	defer s.Lock()
	val, ok := s.data[key]
	if ok {
		return val
	}
	s.data[key] = defaultValue
	return defaultValue
}

func (s *SyncMap[T, K]) Exist(key T) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[key]
	return ok
}

func (s *SyncMap[T, K]) ExistGet(key T) (K, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SyncMap[T, K]) Set(key T, val K) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

func (s *SyncMap[T, K]) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.data)
}

func (s *SyncMap[T, K]) Data() map[T]K {
	return s.data
}

func (s *SyncMap[T, K]) Merge(m map[T]K) {
	s.Lock()
	defer s.Unlock()
	s.data = MergeMap(s.data, m)
}

func (s *SyncMap[T, K]) Clone() *SyncMap[T, K] {
	s.RLock()
	defer s.RUnlock()
	clone := map[T]K{}
	for k, v := range s.data {
		clone[k] = v
	}
	return WrapToSyncMap(clone)
}
