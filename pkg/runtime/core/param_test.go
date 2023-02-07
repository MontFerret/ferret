package core_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestParamsWith(t *testing.T) {
	Convey("Should match", t, func() {
		p := make(map[string]core.Value)
		p["val1"] = values.NewInt(1)
		p["val2"] = values.NewString("test")

		pc := core.ParamsWith(context.Background(), p)
		So(pc, ShouldNotBeNil)

		out, err := core.ParamsFrom(pc)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, p)
	})
}

func TestParamsFrom(t *testing.T) {
	Convey("Should match", t, func() {
		p := make(map[string]core.Value)
		p["val1"] = values.NewInt(1)
		p["val2"] = values.NewString("test")

		_, err := core.ParamsFrom(context.Background())

		So(err, ShouldNotBeNil)

		ctx := context.WithValue(context.Background(), "fail", p)
		pf, err := core.ParamsFrom(ctx)

		So(err, ShouldNotBeNil)
		So(pf, ShouldBeNil)

		ctx = context.WithValue(context.Background(), "params", p)
		pf, err = core.ParamsFrom(ctx)

		So(err, ShouldNotBeNil)
		So(pf, ShouldBeNil)

		ctx = core.ParamsWith(context.Background(), p)
		pf, err = core.ParamsFrom(ctx)

		So(err, ShouldBeNil)
		So(pf, ShouldEqual, p)
	})
}

func TestParamFrom(t *testing.T) {
	Convey("Should match", t, func() {
		p := make(map[string]core.Value)
		p["val1"] = values.NewInt(1)
		p["val2"] = values.NewString("test")

		_, err := core.ParamFrom(context.Background(), "")

		So(err, ShouldNotBeNil)

		ctx := context.WithValue(context.Background(), "fail", p)
		_, err = core.ParamFrom(ctx, "val1")

		So(err, ShouldNotBeNil)

		ctx = core.ParamsWith(context.Background(), p)
		v, err := core.ParamFrom(ctx, "val1")

		So(err, ShouldBeNil)
		So(v, ShouldEqual, values.NewInt(1))
	})
}
