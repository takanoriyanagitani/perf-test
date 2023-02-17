package http_client_perf

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	t.Parallel()

	t.Run("GeneratorStatic", func(t *testing.T) {
		t.Parallel()

		t.Run("strings", func(t *testing.T) {
			t.Parallel()

			var gen GeneratorStatic[[]string] = func() []string {
				return []string{
					"alpha",
					"bravo",
					"charlie",
					"delta",
					"echo",
					"foxtrot",
					"golf",
					"hotel",
					"india",
				}
			}

			var upd MapUpdate = gen.ToMapUpdator("keywords")
			var m MapInput = MapInputNew()
			upd(&m)

			got, found := m.Get("keywords")
			t.Run("found", assertTrue(found))
			keywords, typeOk := got.([]string)
			t.Run("type ok", assertTrue(typeOk))
			t.Run("len check", assertEq(len(keywords), 9))
		})
	})
}

func BenchmarkGenerator(b *testing.B) {
	b.Run("GeneratorStatic", func(b *testing.B) {
		b.Run("strings", func(b *testing.B) {
			var gen GeneratorStatic[[]string] = func() []string {
				return []string{
					"alpha",
					"bravo",
					"charlie",
					"delta",
					"echo",
					"foxtrot",
					"golf",
					"hotel",
					"india",
				}
			}

			var upd MapUpdate = gen.ToMapUpdator("keywords")
			var m MapInput = MapInputNew()

			upd(&m)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				upd(&m)
			}
		})
	})
}
