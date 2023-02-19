package http_client_perf

import (
	"encoding/json"
	"io"
)

type MapInput struct{ m map[string]any }

func MapInputNew() MapInput { return MapInput{m: make(map[string]any)} }

type MapUpdate func(m *MapInput)

func (m *MapInput) Set(key string, val any) { m.m[key] = val }
func (m *MapInput) Get(key string) (val any, found bool) {
	val, found = m.m[key]
	return val, found
}

func (m *MapInput) ToJson(writer io.Writer) error {
	var encoder *json.Encoder = json.NewEncoder(writer)
	return encoder.Encode(m.m)
}

func (m *MapInput) ToMapUpdator(
	updator MapUpdate,
	key string,
) MapUpdate {
	return func(parent *MapInput) {
		updator(m)
		var updated map[string]any = m.m
		parent.Set(key, updated)
	}
}

func (u MapUpdate) Append(others ...MapUpdate) MapUpdate {
	return func(m *MapInput) {
		u(m)
		for _, other := range others {
			other(m)
		}
	}
}

type MapUpdateNew[R any] func(randomSource R) MapUpdate
