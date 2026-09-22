package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func BenchmarkIntegerWidthCalendar(b *testing.B) {
	date := runtime.NewDateTime(time.Date(2024, time.January, 31, 12, 0, 0, 0, time.UTC))
	for _, unit := range []runtime.String{"hour", "day", "month", "year"} {
		b.Run(string(unit), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := datetime.Add(b.Context(), date, runtime.Int(1), unit); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
