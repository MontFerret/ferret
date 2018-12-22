package runtime_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Result struct {
	Value []byte
	Error error
}

func TestProgram(t *testing.T) {
	Convey("Should recover from panic", t, func() {
		c := compiler.New()
		c.RegisterFunction("panic", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			panic("test")
		})

		p := c.MustCompile(`RETURN PANIC()`)

		_, err := p.Run(context.Background())

		So(err, ShouldBeError)
		So(err.Error(), ShouldEqual, "test")
	})

	Convey("Should stop an execution when context is cancelled", t, func() {
		c := compiler.New()
		p := c.MustCompile(`WAIT(1000) RETURN TRUE`)

		out := make(chan Result)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			v, err := p.Run(ctx)

			out <- Result{
				Value: v,
				Error: err,
			}
		}()

		cancel()

		o := <-out

		So(o.Error, ShouldEqual, core.ErrTerminated)
	})
}
