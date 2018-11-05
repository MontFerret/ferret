package runtime_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
}
