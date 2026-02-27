package datetime_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type testCase struct {
	Name      string
	Expected  runtime.Value
	Args      []runtime.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn runtime.Function) {
	Convey(tc.Name, t, func() {
		actual, err := fn(context.Background(), tc.Args...)

		if tc.ShouldErr {
			So(err, ShouldNotBeNil)
		} else {
			So(err, ShouldBeNil)
			So(actual.String(), ShouldEqual, tc.Expected.String())
		}
	})
}

func Fn0(fn runtime.Function0) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		if err := runtime.ValidateArgs(args, 0, 0); err != nil {
			return runtime.None, err
		}

		return fn(ctx)
	}
}

func Fn1(fn runtime.Function1) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		if err := runtime.ValidateArgs(args, 1, 1); err != nil {
			return runtime.None, err
		}

		return fn(ctx, args[0])
	}
}

func Fn2(fn runtime.Function2) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		if err := runtime.ValidateArgs(args, 2, 2); err != nil {
			return runtime.None, err
		}

		return fn(ctx, args[0], args[1])
	}
}

func Fn3(fn runtime.Function3) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		if err := runtime.ValidateArgs(args, 3, 3); err != nil {
			return runtime.None, err
		}

		return fn(ctx, args[0], args[1], args[2])
	}
}

func Fn4(fn runtime.Function4) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		if err := runtime.ValidateArgs(args, 4, 4); err != nil {
			return runtime.None, err
		}

		return fn(ctx, args[0], args[1], args[2], args[3])
	}
}

func mustDefaultLayoutDt(timeString string) runtime.DateTime {
	dt, err := defaultLayoutDt(timeString)

	if err != nil {
		panic(err)
	}

	return dt
}

func mustLayoutDt(layout, value string) runtime.DateTime {
	dt, err := layoutDt(layout, value)

	if err != nil {
		panic(err)
	}

	return dt
}

func defaultLayoutDt(timeString string) (runtime.DateTime, error) {
	return layoutDt(runtime.DefaultTimeLayout, timeString)
}

func layoutDt(layout, value string) (runtime.DateTime, error) {
	t, err := time.Parse(layout, value)

	if err != nil {
		return runtime.DateTime{}, err
	}

	return runtime.NewDateTime(t), nil
}
