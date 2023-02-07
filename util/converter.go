package util

func SetToList[T comparable](m map[T]bool) []T {
	l := []T{}
	for k, v := range m {
		if v {
			l = append(l, k)
		}
	}
	return l
}

func ListToSet[T comparable](l []T) map[T]bool {
	s := map[T]bool{}
	for _, v := range l {
		s[v] = true
	}
	return s
}

func MapToList[T comparable, K any](m map[T]K) []*Pair[T, K] {
	l := []*Pair[T, K]{}
	for k, v := range m {
		l = append(l, NewPair(k, v))
	}
	return l
}
