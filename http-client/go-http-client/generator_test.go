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

type generatorTestBufContainerF32 struct{ internal []float32 }

type generatorTestBufContainerNestedF32iv struct{ internal [][16]float32 }

type generatorTestBufContainerNestedF32s struct{ internal [][]float32 }

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

	b.Run("GeneratorArrayNumBuilderBufferedNew", func(b *testing.B) {
		b.Run("empty", func(b *testing.B) {
			var container generatorTestBufContainerF32 = generatorTestBufContainerF32{
				internal: make([]float32, 0, 8),
			}
			var builder func(
				c *generatorTestBufContainerF32,
			) GeneratorArrayNum[int, float32] = GeneratorArrayNumBuilderBufferedNew(
				func(input int, buf []float32) []float32 { return nil },
				func(c *generatorTestBufContainerF32) { c.internal = c.internal[:0] },
				func(c *generatorTestBufContainerF32) (buf []float32) { return c.internal },
				func(c *generatorTestBufContainerF32, buf []float32) { c.internal = buf },
			)
			var gen GeneratorArrayNum[int, float32] = builder(&container)

			b.ResetTimer()
			var tot int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(i)
				tot += len(generated)
			}
			if 0 != tot {
				b.Fatalf("Must be empty. tot: %v", tot)
			}
		})

		b.Run("max", func(b *testing.B) {
			var container generatorTestBufContainerF32 = generatorTestBufContainerF32{
				internal: make([]float32, 0, 8),
			}
			var builder func(
				c *generatorTestBufContainerF32,
			) GeneratorArrayNum[int, float32] = GeneratorArrayNumBuilderBufferedNew(
				func(input int, buf []float32) []float32 {
					var generated []float32 = buf
					generated = append(generated, 0.0)
					generated = append(generated, 1.0)
					generated = append(generated, 2.0)
					generated = append(generated, 3.0)
					generated = append(generated, 4.0)
					generated = append(generated, 5.0)
					generated = append(generated, 6.0)
					generated = append(generated, 7.0)
					return generated
				},
				func(c *generatorTestBufContainerF32) { c.internal = c.internal[:0] },
				func(c *generatorTestBufContainerF32) (buf []float32) { return c.internal },
				func(c *generatorTestBufContainerF32, buf []float32) { c.internal = buf },
			)
			var gen GeneratorArrayNum[int, float32] = builder(&container)

			b.ResetTimer()
			var tot int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(i)
				tot += len(generated)
			}
			if (b.N * 8) != tot {
				b.Fatalf("Unexpected total: %v", tot)
			}
		})

		b.Run("more items", func(b *testing.B) {
			var container generatorTestBufContainerF32 = generatorTestBufContainerF32{
				internal: make([]float32, 0, 8),
			}
			var builder func(
				c *generatorTestBufContainerF32,
			) GeneratorArrayNum[int, float32] = GeneratorArrayNumBuilderBufferedNew(
				func(input int, buf []float32) []float32 {
					var generated []float32 = buf
					generated = append(generated, 0.0)
					generated = append(generated, 1.0)
					generated = append(generated, 2.0)
					generated = append(generated, 3.0)
					generated = append(generated, 4.0)
					generated = append(generated, 5.0)
					generated = append(generated, 6.0)
					generated = append(generated, 7.0)

					generated = append(generated, 8.0)
					return generated
				},
				func(c *generatorTestBufContainerF32) { c.internal = c.internal[:0] },
				func(c *generatorTestBufContainerF32) (buf []float32) { return c.internal },
				func(c *generatorTestBufContainerF32, buf []float32) { c.internal = buf },
			)
			var gen GeneratorArrayNum[int, float32] = builder(&container)

			b.ResetTimer()
			var tot int = 0
			for i := 0; i < b.N; i++ {
				var generated []float32 = gen(i)
				tot += len(generated)
			}
			if (b.N * 9) != tot {
				b.Fatalf("Unexpected total: %v", tot)
			}
		})
	})

	b.Run("GeneratorArray", func(b *testing.B) {
		b.Run("nested array", func(b *testing.B) {
			b.Run("empty", func(b *testing.B) {
				container := generatorTestBufContainerNestedF32iv{
					internal: make([][16]float32, 256),
				}

				var builder func(c *generatorTestBufContainerNestedF32iv) GeneratorArray[
					int,
					[16]float32,
				] = GeneratorArrayBuilderBufferedNew(
					func(input int, buf [][16]float32) (output [][16]float32) { return },
					func(c *generatorTestBufContainerNestedF32iv) { c.internal = c.internal[:0] },
					func(c *generatorTestBufContainerNestedF32iv) (buf [][16]float32) {
						return c.internal
					},
					func(c *generatorTestBufContainerNestedF32iv, buf [][16]float32) {
						c.internal = buf
					},
				)

				var gen GeneratorArray[int, [16]float32] = builder(&container)

				var tot int = 0
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var generated [][16]float32 = gen(i)
					tot += len(generated)
				}
				if 0 != tot {
					b.Fatalf("Unexpected total: %v", tot)
				}
			})

			b.Run("max", func(b *testing.B) {
				container := generatorTestBufContainerNestedF32iv{
					internal: make([][16]float32, 256),
				}

				var builder func(c *generatorTestBufContainerNestedF32iv) GeneratorArray[
					int,
					[16]float32,
				] = GeneratorArrayBuilderBufferedNew(
					func(input int, buf [][16]float32) (output [][16]float32) {
						var generated [][16]float32 = buf
						var temp [16]float32
						for i := 0; i < cap(buf); i++ {
							for j := 0; j < 16; j++ {
								temp[j] = float32(i) * float32(j)
							}
							generated = append(generated, temp)
						}
						return generated
					},
					func(c *generatorTestBufContainerNestedF32iv) { c.internal = c.internal[:0] },
					func(c *generatorTestBufContainerNestedF32iv) (buf [][16]float32) {
						return c.internal
					},
					func(c *generatorTestBufContainerNestedF32iv, buf [][16]float32) {
						c.internal = buf
					},
				)

				var gen GeneratorArray[int, [16]float32] = builder(&container)

				var tot int = 0
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var generated [][16]float32 = gen(i)
					tot += len(generated)
				}
				if (b.N * 256) != tot {
					b.Fatalf("Unexpected total: %v", tot)
				}
			})

			b.Run("more items", func(b *testing.B) {
				container := generatorTestBufContainerNestedF32iv{
					internal: make([][16]float32, 256),
				}

				const items int = 512

				var builder func(c *generatorTestBufContainerNestedF32iv) GeneratorArray[
					int,
					[16]float32,
				] = GeneratorArrayBuilderBufferedNew(
					func(input int, buf [][16]float32) (output [][16]float32) {
						var generated [][16]float32 = buf
						var temp [16]float32
						for i := 0; i < items; i++ {
							for j := 0; j < 16; j++ {
								temp[j] = float32(i) * float32(j)
							}
							generated = append(generated, temp)
						}
						return generated
					},
					func(c *generatorTestBufContainerNestedF32iv) { c.internal = c.internal[:0] },
					func(c *generatorTestBufContainerNestedF32iv) (buf [][16]float32) {
						return c.internal
					},
					func(c *generatorTestBufContainerNestedF32iv, buf [][16]float32) {
						c.internal = buf
					},
				)

				var gen GeneratorArray[int, [16]float32] = builder(&container)

				var tot int = 0
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var generated [][16]float32 = gen(i)
					tot += len(generated)
				}
				if (b.N * items) != tot {
					b.Fatalf("Unexpected total: %v", tot)
				}
			})
		})

		b.Run("nested slice", func(b *testing.B) {
			b.Run("empty", func(b *testing.B) {
				var width int = 9
				var height int = 256
				var internal [][]float32 = make([][]float32, height)
				var fixedWidth []float32 = make([]float32, height*width)
				for ix := range internal {
					internal[ix], fixedWidth = fixedWidth[:width], fixedWidth[width:]
				}
				container := generatorTestBufContainerNestedF32s{internal}

				var builder func(c *generatorTestBufContainerNestedF32s) GeneratorArray[
					int,
					[]float32,
				] = GeneratorArrayBuilderBufferedNew(
					func(input int, buf [][]float32) (output [][]float32) { return },
					func(c *generatorTestBufContainerNestedF32s) { c.internal = c.internal[:0] },
					func(c *generatorTestBufContainerNestedF32s) (buf [][]float32) {
						return c.internal
					},
					func(c *generatorTestBufContainerNestedF32s, buf [][]float32) {
						c.internal = buf
					},
				)

				var gen GeneratorArray[int, []float32] = builder(&container)

				var tot int = 0
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var generated [][]float32 = gen(i)
					tot += len(generated)
				}
				if 0 != tot {
					b.Fatalf("Unexpected total: %v", tot)
				}
			})

			b.Run("max", func(b *testing.B) {
				var width int = 9
				var height int = 256
				var internal [][]float32 = make([][]float32, height)
				var fixedWidth []float32 = make([]float32, height*width)
				for ix := range internal {
					internal[ix], fixedWidth = fixedWidth[:width], fixedWidth[width:]
				}
				container := generatorTestBufContainerNestedF32s{internal}

				var builder func(c *generatorTestBufContainerNestedF32s) GeneratorArray[
					int,
					[]float32,
				] = GeneratorArrayBuilderBufferedNew(
					func(input int, buf [][]float32) (output [][]float32) {
						output = buf
						for i := 0; i < len(buf); i++ {
							var row []float32 = buf[i]
							for j := 0; j < len(row); j++ {
								row[j] = float32(i) * float32(j)
							}
						}
						return
					},
					func(c *generatorTestBufContainerNestedF32s) {}, // nop
					func(c *generatorTestBufContainerNestedF32s) (buf [][]float32) {
						return c.internal
					},
					func(c *generatorTestBufContainerNestedF32s, buf [][]float32) {
						c.internal = buf
					},
				)

				var gen GeneratorArray[int, []float32] = builder(&container)

				var tot int = 0
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var generated [][]float32 = gen(i)
					tot += len(generated)
				}
				if (b.N * height) != tot {
					b.Fatalf("Unexpected total: %v", tot)
				}
			})
		})
	})
}
