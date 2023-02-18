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

type GeneratorArrayNum[I any, N constraints.Integer | constraints.Float] func(input I) []N

func GeneratorArrayNumBuilderNew[I any, N constraints.Integer | constraints.Float](
	gen func(input I, buf []N) []N,
) func(buf []N) GeneratorArrayNum[I, N] {
	return func(buf []N) GeneratorArrayNum[I, N] {
		var b []N = buf
		updateBuf := func(neo []N) { b = neo }
		return func(input I) []N {
			b = b[:0]
			var neo []N = gen(input, b)
			if cap(neo) != cap(buf) {
				updateBuf(neo)
			}
			return neo
		}
	}
}

func GeneratorArrayNumBuilderBufferedNew[B any, I any, N constraints.Integer | constraints.Float](
	gen func(input I, buf []N) []N,
	resetBuf func(container B),
	getBuf func(container B) (buf []N),
	updateBuf func(container B, buf []N),
) func(container B) GeneratorArrayNum[I, N] {
	return func(container B) GeneratorArrayNum[I, N] {
		return func(input I) []N {
			resetBuf(container)
			var buf []N = getBuf(container)
			var generated []N = gen(input, buf)
			updateBuf(container, generated)
			return generated
		}
	}
}

type GeneratorArray[I, N any] func(input I) (output []N)

func (g GeneratorArray[I, N]) ToMapUpdator(
	inputGen func() (input I),
	key string,
) MapUpdate {
	return func(i *MapInput) {
		var input I = inputGen()
		var output []N = g(input)
		i.Set(key, output)
	}
}

func GeneratorArrayBuilderBufferedNew[B, I, N any](
	gen func(input I, buf []N) (output []N),
	resetBuf func(container B),
	getBuf func(container B) (buf []N),
	updateBuf func(container B, buf []N),
) func(container B) GeneratorArray[I, N] {
	return func(container B) GeneratorArray[I, N] {
		return func(input I) (output []N) {
			resetBuf(container)
			var buf []N = getBuf(container)
			var generated []N = gen(input, buf)
			updateBuf(container, generated)
			return generated
		}
	}
}
