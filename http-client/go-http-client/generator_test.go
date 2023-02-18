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

	t.Run("GeneratorArrayNum", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(input float64, buf []float32) []float32 { return nil },
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])
			var generated []float32 = gen(0.0)
			t.Run("no items", assertEmpty(generated))
		})

		t.Run("single item", func(t *testing.T) {
			t.Parallel()

			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(input float64, buf []float32) []float32 {
					buf = append(buf, float32(input))
					return buf
				},
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])
			var generated []float32 = gen(0.0)
			t.Run("single item", assertEq(len(generated), 1))
			t.Run("same", assertEq(generated[0], 0.0))
		})

		t.Run("max item", func(t *testing.T) {
			t.Parallel()

			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(input float64, buf []float32) []float32 {
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))

					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					return buf
				},
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])
			var generated []float32 = gen(0.0)
			t.Run("len check", assertEq(len(generated), 8))
		})

		t.Run("more items", func(t *testing.T) {
			t.Parallel()

			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(input float64, buf []float32) []float32 {
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))

					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))
					buf = append(buf, float32(input))

					buf = append(buf, float32(input))
					return buf
				},
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])
			var generated []float32 = gen(0.0)
			t.Run("len check", assertEq(len(generated), 9))
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

	b.Run("GeneratorArrayNum", func(b *testing.B) {
		b.Run("empty", func(b *testing.B) {
			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(input float64, buf []float32) []float32 { return nil },
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])

			b.ResetTimer()
			var tot int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(0.0)
				tot += len(generated)
			}
			if 0 != tot {
				b.Fatalf("Must be empty")
			}
		})

		b.Run("max", func(b *testing.B) {
			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(_ float64, buf []float32) []float32 {
					buf = append(buf, 0.0)
					buf = append(buf, 1.0)
					buf = append(buf, 2.0)
					buf = append(buf, 3.0)
					buf = append(buf, 4.0)
					buf = append(buf, 5.0)
					buf = append(buf, 6.0)
					buf = append(buf, 7.0)
					return buf
				},
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])

			b.ResetTimer()
			var tot int = 0
			var caps int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(0.0)
				tot += len(generated)
				caps += cap(generated)
			}
			if (b.N * 8) != tot {
				b.Fatalf("Unexpected total: %v", tot)
			}
			if (b.N * 8) != caps {
				b.Fatalf("Unexpected capacity: %v", caps)
			}
		})

		b.Run("more items", func(b *testing.B) {
			var buf [8]float32
			var builder func(
				buf []float32,
			) GeneratorArrayNum[float64, float32] = GeneratorArrayNumBuilderNew(
				func(_ float64, buf []float32) []float32 {
					buf = append(buf, 0.0)
					buf = append(buf, 1.0)
					buf = append(buf, 2.0)
					buf = append(buf, 3.0)
					buf = append(buf, 4.0)
					buf = append(buf, 5.0)
					buf = append(buf, 6.0)
					buf = append(buf, 7.0)

					buf = append(buf, 8.0)
					return buf
				},
			)
			var gen GeneratorArrayNum[float64, float32] = builder(buf[:])

			b.ResetTimer()
			var tot int = 0
			var caps int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(0.0)
				tot += len(generated)
				caps += cap(generated)
			}
			if (b.N * 9) != tot {
				b.Fatalf("Unexpected total: %v", tot)
			}
			if (b.N * 8) > caps {
				b.Fatalf("Unexpected capacity: %v", caps)
			}
		})
	})
}
