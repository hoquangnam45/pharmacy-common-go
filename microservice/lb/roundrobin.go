package lb

type RoundRobinLB[T comparable] struct {
	*baseLB[T]
	idx int
}

func NewRoundRobinLB[T comparable]() *RoundRobinLB[T] {
	return &RoundRobinLB[T]{
		baseLB: &baseLB[T]{},
	}
}

func (l *RoundRobinLB[T]) Get() (T, error) {
	err := l.Check()
	var noop T
	if err != nil {
		return noop, err
	}
	ret := l.elements[l.idx]
	l.idx = (l.idx + 1) % len(l.elements)
	return ret, nil
}

func (l *RoundRobinLB[T]) Remove(val T) {
	if idx, ok := l.baseLB.Remove(val); ok {
		if l.idx > idx {
			l.idx -= 1
		}
	}
}
