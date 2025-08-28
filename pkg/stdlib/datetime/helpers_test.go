package datetime_test

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type testCase struct {
	Name      string
	Expected  runtime.Value
	Args      []runtime.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn runtime.Function) {
	t.Run(tc.Name, func(t *testing.T) {
		actual, err := fn(context.Background(), tc.Args...)

		if tc.ShouldErr {
			if err == nil {
				t.Errorf("expected error but got none")
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if actual.String() != tc.Expected.String() {
				t.Errorf("expected %v but got %v", tc.Expected, actual)
			}
		}
	})
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
