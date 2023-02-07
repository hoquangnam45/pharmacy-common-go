package lb

func createPositionMap[T comparable](elements []T) map[T]int {
	positionMap := map[T]int{}
	for i, v := range elements {
		positionMap[v] = i
	}
	return positionMap
}
