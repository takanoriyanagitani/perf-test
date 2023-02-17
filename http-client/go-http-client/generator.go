package http_client_perf

import (
	"golang.org/x/exp/constraints"
)

type GeneratorLite[I, O constraints.Integer | constraints.Float] func(i I) O

func (g GeneratorLite[I, O]) UpdateMap(i I, m *MapInput, key string) { m.Set(key, g(i)) }

func (g GeneratorLite[I, O]) ToMapUpdator(
	inputGenerator func() I,
	key string,
) MapUpdate {
	return func(m *MapInput) {
		var i I = inputGenerator()
		g.UpdateMap(i, m, key)
	}
}

type GeneratorStatic[T any] func() T

func GeneratorStaticNew[T any](t T) GeneratorStatic[T] { return func() T { return t } }

func (s GeneratorStatic[T]) ToMapUpdator(key string) MapUpdate {
	return func(m *MapInput) { m.Set(key, s()) }
}
