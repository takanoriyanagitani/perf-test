package go_http_client_test

import (
	"fmt"

	"strings"

	perf "github.com/takanoriyanagitani/perf-test/http-client/go-http-client/v3"
)

type ExampleBufContainerFixed3 struct {
	buf [][8]float32
}

func (b *ExampleBufContainerFixed3) Reset()                  { b.buf = b.buf[:0] }
func (b *ExampleBufContainerFixed3) Get() [][8]float32       { return b.buf }
func (b *ExampleBufContainerFixed3) Update(neo [][8]float32) { b.buf = neo }

func ExampleMapInput() {
	var builder func(
		container *ExampleBufContainerFixed3,
	) perf.GeneratorArray[int, [8]float32] = perf.GeneratorArrayBuilderBufferedNew(
		func(input int, buf [][8]float32) (output [][8]float32) {
			output = buf
			output = append(output, [8]float32{
				0, 1, 2, 3, 4, 5, 6, 7,
			})
			output = append(output, [8]float32{
				8, 9, 0, 1, 2, 3, 4, 5,
			})
			return
		},
		func(c *ExampleBufContainerFixed3) { c.Reset() },
		func(c *ExampleBufContainerFixed3) (buf [][8]float32) { return c.Get() },
		func(c *ExampleBufContainerFixed3, neo [][8]float32) { c.Update(neo) },
	)

	var container ExampleBufContainerFixed3 = ExampleBufContainerFixed3{
		buf: make([][8]float32, 0, 256),
	}
	var gen perf.GeneratorArray[int, [8]float32] = builder(&container)
	var u perf.MapUpdate = gen.ToMapUpdator(
		func() (input int) { return 634 },
		"slice_of_array",
	)
	var i perf.MapInput = perf.MapInputNew()
	u(&i)

	var serialized strings.Builder
	e := i.ToJson(&serialized)
	if nil != e {
		panic(e)
	}

	fmt.Printf("%s", serialized.String())
	// Output: {"slice_of_array":[[0,1,2,3,4,5,6,7],[8,9,0,1,2,3,4,5]]}
}
