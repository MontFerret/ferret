package datetime_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestCalendarShiftBoundaryDSTGap(t *testing.T) {
	const maxEpochSeconds int64 = math.MaxInt64 - 62_135_596_800
	for _, tc := range []struct {
		name       string
		transition int64
		before     int32
		after      int32
		source     int64
		want       int64
	}{
		{name: "wall_above_upper", transition: maxEpochSeconds - 1800, before: 3600, after: 7200, source: maxEpochSeconds - 86400, want: maxEpochSeconds},
		{name: "provisional_above_upper", transition: maxEpochSeconds + 1800, before: -3600, after: 0, source: maxEpochSeconds + 2700 - 86400, want: maxEpochSeconds - 900},
	} {
		t.Run(tc.name, func(t *testing.T) {
			location := calendarTransitionLocation(t, tc.transition, tc.before, tc.after)
			base := runtime.NewDateTime(time.Unix(tc.source, 123).In(location))
			got, err := datetime.Add(t.Context(), base, runtime.Int(1), runtime.String("day"))
			if math.MaxInt == math.MaxInt32 {
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
					t.Fatalf("got %v, %v; want None and argument-1 range error", got, err)
				}

				return
			}

			date, ok := got.(runtime.DateTime)
			if err != nil || !ok {
				t.Fatalf("got %T, %v; want DateTime", got, err)
			}

			if date.Unix() != tc.want || date.Nanosecond() != 123 || date.Location() != location || !date.Equal(base.AddDate(0, 0, 1)) {
				t.Fatalf("got (%d, %d, %p); want (%d, 123, %p) and Go AddDate normalization", date.Unix(), date.Nanosecond(), date.Location(), tc.want, location)
			}
		})
	}
}

func TestCalendarShiftEpochBoundaries(t *testing.T) {
	const (
		maxEpochSeconds int64 = math.MaxInt64 - 62_135_596_800
		// Unlike MinInt64 itself, this epoch has calendar fields that Go can
		// extract without wrapping on a 64-bit host.
		lowerCalendarSeconds int64 = math.MinInt64 + 1<<37
	)

	for _, offset := range []int{-172800, -3600, 0, 3600} {
		location := time.FixedZone(fmt.Sprintf("UTC%+d", offset), offset)
		for _, tc := range []struct {
			name      string
			seconds   int64
			days      runtime.Int
			want      int64
			wantRange bool
		}{
			{name: "upper_zero", seconds: maxEpochSeconds, want: maxEpochSeconds},
			{name: "into_upper", seconds: maxEpochSeconds - 86400, days: 1, want: maxEpochSeconds},
			{name: "from_upper", seconds: maxEpochSeconds, days: -1, want: maxEpochSeconds - 86400},
			{name: "above_upper", seconds: maxEpochSeconds, days: 1, wantRange: true},
			{name: "unix_wrap", seconds: maxEpochSeconds, days: 62_135_596_800/86400 + 1, wantRange: true},
			{name: "lower_wrapped_source", seconds: math.MinInt64, wantRange: true},
			{name: "near_lower_wrapped_source", seconds: math.MinInt64 + 86400, wantRange: true},
			{name: "lower_calendar_zero", seconds: lowerCalendarSeconds, want: lowerCalendarSeconds},
			{name: "lower_calendar_inward", seconds: lowerCalendarSeconds, days: 1, want: lowerCalendarSeconds + 86400},
			{name: "lower_calendar_outward", seconds: lowerCalendarSeconds, days: -1, want: lowerCalendarSeconds - 86400},
			{name: "below_lower", seconds: lowerCalendarSeconds, days: -(1<<37)/86400 - 1, wantRange: true},
		} {
			units := []runtime.String{"day"}
			if tc.days == 0 {
				units = append(units, "week", "month", "year")
			}

			for _, unit := range units {
				for _, subtract := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/%s/%s/subtract=%v", location, tc.name, unit, subtract), func(t *testing.T) {
						function := datetime.Add
						amount := tc.days
						if subtract {
							function = datetime.Subtract
							amount = -amount
						}

						base := runtime.NewDateTime(time.Unix(tc.seconds, 123).In(location))
						got, err := function(t.Context(), base, amount, unit)
						if tc.wantRange || math.MaxInt == math.MaxInt32 {
							position, attributed, _ := runtime.InvalidArgumentDetails(err)
							if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
								t.Fatalf("got %v, %v; want None and argument-1 range error", got, err)
							}

							return
						}

						date, ok := got.(runtime.DateTime)
						if err != nil || !ok {
							t.Fatalf("got %T, %v; want DateTime", got, err)
						}

						if date.Unix() != tc.want || date.Nanosecond() != 123 || date.Location() != location {
							t.Fatalf("got (%d, %d, %p); want (%d, 123, %p)", date.Unix(), date.Nanosecond(), date.Location(), tc.want, location)
						}
					})
				}
			}
		}
	}
}

func TestDateTimeShiftIntegerWidth(t *testing.T) {
	base := time.Date(2024, time.February, 28, 12, 30, 15, 123, time.UTC)
	for _, amount := range []runtime.Int{math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1 << 32} {
		for _, subtract := range []bool{false, true} {
			function := datetime.Add
			delta := time.Duration(amount) * time.Millisecond
			if subtract {
				function = datetime.Subtract
				delta = -delta
			}

			got, err := function(t.Context(), runtime.NewDateTime(base), amount, runtime.String("millisecond"))
			date, ok := got.(runtime.DateTime)
			if err != nil || !ok || !date.Equal(base.Add(delta)) {
				t.Errorf("shift %d, subtract=%v: %v, %v", amount, subtract, got, err)
			}
		}
	}
}

func TestDateTimeShiftRangeErrors(t *testing.T) {
	base := runtime.NewDateTime(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	for _, amount := range []runtime.Int{math.MinInt64, math.MaxInt64, 1 << 50, math.MaxInt64 / 7} {
		for _, unit := range []runtime.String{"millisecond", "second", "minute", "hour", "day", "week", "month", "year"} {
			for _, function := range []runtime.Function3{datetime.Add, datetime.Subtract} {
				got, err := function(t.Context(), base, amount, unit)
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
					t.Errorf("shift(%d, %s) = %v, %v; want argument-1 range error", amount, unit, got, err)
				}
			}
		}
	}
}

func TestCalendarShiftNativeBoundary(t *testing.T) {
	base := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	amount := runtime.Int(math.MaxInt32) + 1
	got, err := datetime.Add(t.Context(), runtime.NewDateTime(base), amount, runtime.String("day"))
	if math.MaxInt == math.MaxInt32 {
		if !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("oversized calendar delta = %v, %v", got, err)
		}
	} else {
		date, ok := got.(runtime.DateTime)
		if err != nil || !ok || !date.Equal(base.AddDate(0, 0, int(amount))) {
			t.Fatalf("representable calendar delta = %v, %v", got, err)
		}
	}

	// The delta itself fits, but AddDate's year addition must also fit.
	_, err = datetime.Add(t.Context(), runtime.NewDateTime(base), runtime.Int(math.MaxInt), runtime.String("year"))
	if !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("calendar component overflow = %v", err)
	}
}

func TestCalendarShiftPreservesNormalization(t *testing.T) {
	for _, zone := range []string{"UTC", "America/New_York", "Pacific/Apia"} {
		location, err := time.LoadLocation(zone)
		if err != nil {
			t.Fatal(err)
		}

		for _, base := range []time.Time{
			time.Date(-10001, time.January, 31, 12, 0, 0, 123, location),
			time.Date(1000000, time.February, 29, 12, 0, 0, 123, location),
			time.Date(2024, time.January, 31, 12, 0, 0, 123, location),
			time.Date(2024, time.February, 29, 12, 0, 0, 123, location),
			time.Date(2024, time.March, 9, 2, 30, 0, 123, location),
			time.Date(2024, time.November, 2, 1, 30, 0, 123, location),
			time.Date(2011, time.December, 29, 12, 0, 0, 123, location),
		} {
			for _, amount := range []int{-3660001, -401, -13, -1, 0, 1, 13, 401, 3660001} {
				for _, unit := range []runtime.String{"day", "week", "month", "year"} {
					var years, months, days int
					switch unit {
					case "day":
						days = amount
					case "week":
						days = amount * 7
					case "month":
						months = amount
					case "year":
						years = amount
					}

					for _, subtract := range []bool{false, true} {
						function := datetime.Add
						direction := 1
						if subtract {
							function = datetime.Subtract
							direction = -1
						}

						want := base.AddDate(years*direction, months*direction, days*direction)
						got, err := function(t.Context(), runtime.NewDateTime(base), runtime.Int(amount), unit)
						date, ok := got.(runtime.DateTime)
						if err != nil || !ok || !date.Equal(want) || date.Location() != location {
							t.Fatalf("%s %v shift %d %s subtract=%v = %v, %v; want %v", zone, base, amount, unit, subtract, got, err, want)
						}
					}
				}
			}
		}
	}
}

func TestCalendarShiftRejectsNarrowedSourceYear(t *testing.T) {
	for _, seconds := range []int64{-(1 << 60), 1 << 60} {
		base := time.Unix(seconds, 123).UTC()
		for _, amount := range []runtime.Int{0, 1} {
			got, err := datetime.Add(t.Context(), runtime.NewDateTime(base), amount, runtime.String("day"))
			if math.MaxInt == math.MaxInt32 {
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
					t.Fatalf("narrowed source year: %v, %v; want argument-1 range error", got, err)
				}
			} else {
				date, ok := got.(runtime.DateTime)
				if err != nil || !ok || !date.Equal(base.AddDate(0, 0, int(amount))) {
					t.Fatalf("representable source year: %v, %v", got, err)
				}
			}
		}
	}
}

// A TZif v2 location with one precisely placed transition, independent of the
// host's zone database. The second block supplies its signed 64-bit epoch.
func calendarTransitionLocation(t *testing.T, transition int64, before, after int32) *time.Location {
	t.Helper()

	var data bytes.Buffer
	for _, wide := range []bool{false, true} {
		data.WriteString("TZif2")
		data.Write(make([]byte, 15))
		for _, count := range []uint32{0, 0, 0, 1, 2, 4} {
			data.Write(binary.BigEndian.AppendUint32(nil, count))
		}

		if wide {
			data.Write(binary.BigEndian.AppendUint64(nil, uint64(transition)))
		} else {
			data.Write(make([]byte, 4))
		}

		data.WriteByte(1)
		data.Write(binary.BigEndian.AppendUint32(nil, uint32(before)))
		data.Write([]byte{0, 0})
		data.Write(binary.BigEndian.AppendUint32(nil, uint32(after)))
		data.Write([]byte{1, 2})
		data.WriteString("A\x00B\x00")
	}

	location, err := time.LoadLocationFromTZData("boundary-gap", data.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	return location
}
