package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestNow(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When too many arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, Fn0(datetime.Now))
	}
}

func TestNowBounds(t *testing.T) {
	functions := datetimeFunctions(t)
	for _, name := range []string{"now", "datetime::now"} {
		fn, ok := functions.A0().Get(name)
		if !ok {
			t.Fatal(name)
		}

		before := time.Now()
		got, err := fn(t.Context())
		after := time.Now()
		date, ok := got.(runtime.DateTime)
		if err != nil || !ok || date.Before(before) || date.After(after) {
			t.Fatalf("%s = %v, %v; outside clock bounds", name, got, err)
		}
	}
}
