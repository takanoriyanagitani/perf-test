package http_client_perf

type MapInput struct{ m map[string]any }

func MapInputNew() MapInput { return MapInput{m: make(map[string]any)} }

func (m *MapInput) Set(key string, val any) { m.m[key] = val }
func (m *MapInput) Get(key string) (val any, found bool) {
	val, found = m.m[key]
	return val, found
}

type MapUpdate func(m *MapInput)

func (u MapUpdate) Append(others ...MapUpdate) MapUpdate {
	return func(m *MapInput) {
		u(m)
		for _, other := range others {
			other(m)
		}
	}
}

type MapUpdateNew[R any] func(randomSource R) MapUpdate
