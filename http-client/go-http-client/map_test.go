package go_http_client

import (
	"testing"

	"bytes"
	"encoding/json"
	"strings"
)

func TestMap(t *testing.T) {
	t.Parallel()

	t.Run("MapInput", func(t *testing.T) {
		t.Parallel()

		t.Run("ToJson", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var m MapInput = MapInputNew()
				var buf strings.Builder
				e := m.ToJson(&buf)
				t.Run("no error", assertNil(e))
				t.Run("empty map", assertEq(buf.String(), "{}\n"))
			})

			t.Run("single flat item", func(t *testing.T) {
				t.Parallel()

				var m MapInput = MapInputNew()
				var genLite GeneratorLite[int, float32] = func(input int) (output float32) {
					return 0.125
				}
				var mapUpdator MapUpdate = genLite.ToMapUpdator(
					func() (input int) { return 634 },
					"static_float",
				)
				mapUpdator(&m)

				var buf strings.Builder
				e := m.ToJson(&buf)
				t.Run("no error", assertNil(e))
				t.Run("a map with a single flat item", assertEq(
					buf.String(),
					`{"static_float":0.125}`+"\n",
				))
			})

			t.Run("multi flat items", func(t *testing.T) {
				t.Parallel()

				var m MapInput = MapInputNew()
				var genLite GeneratorLite[int, float32] = func(input int) (output float32) {
					return 0.125
				}
				var mapUpdator MapUpdate = genLite.ToMapUpdator(
					func() (input int) { return 634 },
					"static_float",
				)

				var genInt GeneratorLite[float64, int32] = func(input float64) (output int32) {
					return 3776
				}
				var allUpdator MapUpdate = mapUpdator.Append(
					genInt.ToMapUpdator(
						func() (input float64) { return 2.99792458 },
						"static_int",
					),
				)
				allUpdator(&m)

				var buf bytes.Buffer
				e := m.ToJson(&buf)
				t.Run("no error", assertNil(e))

				var got map[string]any
				e = json.Unmarshal(buf.Bytes(), &got)
				t.Run("no deser error", assertNil(e))

				_staticFloat, found := got["static_float"]
				t.Run("float like found", assertTrue(found))
				staticFloat, ok := _staticFloat.(float64)
				t.Run("float found", assertTrue(ok))
				t.Run("same float", assertEq(staticFloat, 0.125))

				_staticInt, found := got["static_int"]
				t.Run("int like item found", assertTrue(found))
				staticInt, ok := _staticInt.(float64)
				t.Run("int like found", assertTrue(ok))
				t.Run("same int", assertEq(int32(staticInt), 3776))
			})

			t.Run("single array", func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					var dummyBufContainer uint8 = 0
					var back [8]float32
					var buf []float32 = back[:]

					var builder func(
						_dummy uint8,
					) GeneratorArray[int64, float32] = GeneratorArrayBuilderBufferedNew(
						func(_input int64, _buf []float32) (output []float32) { return },
						func(_dummy uint8) { buf = buf[:0] },
						func(_dummy uint8) []float32 { return buf },
						func(_dummy uint8, neo []float32) { buf = neo },
					)

					var gen GeneratorArray[int64, float32] = builder(dummyBufContainer)
					var u MapUpdate = gen.ToMapUpdator(
						func() (input int64) { return 299792458 },
						"empty_array",
					)

					var i MapInput = MapInputNew()
					u(&i)

					var serialized bytes.Buffer
					e := i.ToJson(&serialized)
					t.Run("no error", assertNil(e))

					typed := struct {
						EmptyArray []float64 `json:"empty_array"`
					}{}
					e = json.Unmarshal(serialized.Bytes(), &typed)
					t.Run("no parse error", assertNil(e))
					t.Run("empty array", assertEmpty(typed.EmptyArray))
				})

				t.Run("single item", func(t *testing.T) {
					t.Parallel()

					var dummyBufContainer uint8 = 0
					var back [8]float32
					var buf []float32 = back[:]

					var builder func(
						_dummy uint8,
					) GeneratorArray[int64, float32] = GeneratorArrayBuilderBufferedNew(
						func(_input int64, buf []float32) (output []float32) {
							output = buf
							output = append(output, 0.75)
							return
						},
						func(_dummy uint8) { buf = buf[:0] },
						func(_dummy uint8) []float32 { return buf },
						func(_dummy uint8, neo []float32) { buf = neo },
					)

					var gen GeneratorArray[int64, float32] = builder(dummyBufContainer)
					var u MapUpdate = gen.ToMapUpdator(
						func() (input int64) { return 299792458 },
						"single_item",
					)

					var i MapInput = MapInputNew()
					u(&i)

					var serialized bytes.Buffer
					e := i.ToJson(&serialized)
					t.Run("no error", assertNil(e))

					typed := struct {
						SingleItem []float64 `json:"single_item"`
					}{}
					e = json.Unmarshal(serialized.Bytes(), &typed)
					t.Run("no parse error", assertNil(e))
					t.Run("single item", assertEq(len(typed.SingleItem), 1))
					var item float64 = typed.SingleItem[0]
					t.Run("same item", assertEq(item, 0.75))
				})
			})
		})
	})
}
