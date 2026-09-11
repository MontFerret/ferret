package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func BenchmarkDateTime(b *testing.B) {
	ns := runtime.NewLibrary()
	datetime.RegisterLib(ns)
	functions, err := ns.Build()
	if err != nil {
		b.Fatal(err)
	}

	start := runtime.NewDateTime(time.Date(2024, time.February, 28, 12, 30, 0, 0, time.UTC))
	end := runtime.NewDateTime(start.Add(90 * time.Minute))
	for _, name := range []string{"date", "date_year", "date_days_in_month"} {
		b.Run(name, func(b *testing.B) {
			fn, ok := functions.A1().Get(name)
			if !ok {
				b.Fatal(name)
			}

			var arg runtime.Value = start
			if name == "date" {
				arg = runtime.String("2024-02-28T12:30:00Z")
			}

			for b.Loop() {
				if _, err := fn(b.Context(), arg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for _, name := range []string{"date_add", "date_diff", "date_compare"} {
		b.Run(name, func(b *testing.B) {
			fn, ok := functions.A3().Get(name)
			if !ok {
				b.Fatal(name)
			}

			var right runtime.Value = end
			if name == "date_add" {
				right = runtime.Int(2)
			}

			for b.Loop() {
				if _, err := fn(b.Context(), start, right, runtime.String("hour")); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDateTimeSame(b *testing.B) {
	base := time.Date(2024, time.June, 12, 10, 20, 30, 123456789, time.UTC)
	var left runtime.Value = runtime.NewDateTime(base)
	var right runtime.Value = runtime.NewDateTime(base.Add(time.Nanosecond))
	for _, unit := range []runtime.String{"year", "week", "millisecond"} {
		b.Run(string(unit), func(b *testing.B) {
			var argument runtime.Value = unit
			for b.Loop() {
				value, err := datetime.Same(b.Context(), left, right, argument)
				if err != nil || value != runtime.True {
					b.Fatalf("Same = %v, %v", value, err)
				}
			}
		})
	}
}
