package util

import (
	"errors"
	"math/rand"
)

func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func RandomMapEntry[T comparable, K any](m map[T]K) (*Pair[T, K], error) {
	var noopT T
	var noopK K
	if len(m) == 0 {
		return NewPair(noopT, noopK), errors.New("empty map")
	}
	list := MapToList(m)
	return list[rand.Intn(len(m))], nil
}

func MergeMap[T comparable, K any](m1 map[T]K, m2 map[T]K) map[T]K {
	ret := map[T]K{}
	for k, v := range m1 {
		ret[k] = v
	}
	for k, v := range m2 {
		ret[k] = v
	}
	return ret
}

func MergeList[T any](l1 []T, l2 []T) []T {
	ret := []T{}
	ret = append(ret, l1...)
	return append(ret, l2...)
}
