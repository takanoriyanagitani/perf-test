package go_http_client

import (
	"time"

	"golang.org/x/exp/constraints"
)

type TimeGeneratorLite[T constraints.Integer | constraints.Float] func(t time.Time) T

func (g TimeGeneratorLite[T]) UpdateMap(t time.Time, m *MapInput, key string) { m.Set(key, g(t)) }

func (g TimeGeneratorLite[T]) ToMapUpdator(
	inputGenerator func() time.Time,
	key string,
) MapUpdate {
	return func(m *MapInput) {
		var t time.Time = inputGenerator()
		g.UpdateMap(t, m, key)
	}
}

func (g TimeGeneratorLite[T]) ToMapUpdatorDefault(key string) MapUpdate {
	return g.ToMapUpdator(time.Now, key)
}
