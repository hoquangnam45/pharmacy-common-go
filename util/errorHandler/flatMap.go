package errorHandler

func FlatMap[T any, K any](m *MaybeError[T], f func(T) *MaybeError[K]) *MaybeError[K] {
	return FactoryM(func() (K, error) {
		val, err := m.Eval()
		var noop K
		if err != nil {
			return noop, err
		}
		newVal, newErr := f(val).Eval()
		return newVal, newErr
	})
}

func FlatMap2[A any, B any, C any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C]) *MaybeError[C] {
	return FlatMap(FlatMap(m, a), b)
}

func FlatMap3[A any, B any, C any, D any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D]) *MaybeError[D] {
	return FlatMap(FlatMap2(m, a, b), c)
}

func FlatMap4[A any, B any, C any, D any, E any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E]) *MaybeError[E] {
	return FlatMap(FlatMap3(m, a, b, c), d)
}

func FlatMap5[A any, B any, C any, D any, E any, F any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F]) *MaybeError[F] {
	return FlatMap(FlatMap4(m, a, b, c, d), e)
}

func FlatMap6[A any, B any, C any, D any, E any, F any, G any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G]) *MaybeError[G] {
	return FlatMap(FlatMap5(m, a, b, c, d, e), f)
}

func FlatMap7[A any, B any, C any, D any, E any, F any, G any, H any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H]) *MaybeError[H] {
	return FlatMap(FlatMap6(m, a, b, c, d, e, f), g)
}

func FlatMap8[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H],
	h func(H) *MaybeError[I]) *MaybeError[I] {
	return FlatMap(FlatMap7(m, a, b, c, d, e, f, g), h)
}

func FlatMap9[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[G],
	g func(G) *MaybeError[H],
	h func(H) *MaybeError[I],
	i func(I) *MaybeError[J]) *MaybeError[J] {
	return FlatMap(FlatMap8(m, a, b, c, d, e, f, g, h), i)
}

func FlatMap10[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any](
	m *MaybeError[A],
	a func(A) *MaybeError[B],
	b func(B) *MaybeError[C],
	c func(C) *MaybeError[D],
	d func(D) *MaybeError[E],
	e func(E) *MaybeError[F],
	f func(F) *MaybeError[J],
	g func(J) *MaybeError[H],
	h func(H) *MaybeError[I],
	i func(I) *MaybeError[J],
	j func(J) *MaybeError[K]) *MaybeError[K] {
	return FlatMap(FlatMap9(m, a, b, c, d, e, f, g, h, i), j)
}
