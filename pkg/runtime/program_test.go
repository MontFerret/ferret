package runtime_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
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

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := p.Run(ctx)

		So(err, ShouldEqual, core.ErrTerminated)
	})
}
