package datetime

import (
	"context"
	"math"
	"math/big"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func addCalendar(ctx context.Context, date time.Time, amount runtime.Int, u unit) (time.Time, error) {
	var years, months, days runtime.Int
	switch u {
	case day:
		days = amount
	case week:
		if amount < math.MinInt64/7 || amount > math.MaxInt64/7 {
			return time.Time{}, calendarRangeError()
		}

		days = amount * 7
	case month:
		months = amount
	case year:
		years = amount
	default:
		return time.Time{}, runtime.Errorf(runtime.ErrUnexpected, "unsupported datetime unit %d", u)
	}

	nativeYears, yearsOK := runtime.ToNativeInt(years)
	nativeMonths, monthsOK := runtime.ToNativeInt(months)
	nativeDays, daysOK := runtime.ToNativeInt(days)
	if !yearsOK || !monthsOK || !daysOK {
		return time.Time{}, calendarRangeError()
	}

	y, m, d := date.Date()
	// Date returns a native year. For epochs outside this inner interval,
	// verify that extracting the calendar fields did not narrow it on 32-bit
	// hosts, even when the requested shift is zero.
	if seconds := date.Unix(); seconds < -(1<<55) || seconds > 1<<55 {
		wall, err := calendarWallTime(ctx, date, y, m, d)
		if err != nil {
			return time.Time{}, err
		}

		_, offset := date.Zone()
		restored, err := calendarUnixSeconds(ctx, wall.Unix(), offset)
		if err != nil || restored != seconds {
			return time.Time{}, calendarRangeError()
		}
	}

	y, yearOK := calendarComponent(y, nativeYears)
	mn, monthOK := calendarComponent(int(m), nativeMonths)
	d, dayOK := calendarComponent(d, nativeDays)
	if !yearOK || !monthOK || !dayOK || mn == math.MinInt || d == math.MinInt {
		return time.Time{}, calendarRangeError()
	}

	// Date normalizes months before computing the day. Check that its native
	// year arithmetic is safe as well as AddDate's component additions.
	monthIndex := mn - 1
	carry, monthIndex := monthIndex/12, monthIndex%12
	if monthIndex < 0 {
		carry--
		monthIndex += 12
	}

	y, yearOK = calendarComponent(y, carry)
	if !yearOK {
		return time.Time{}, calendarRangeError()
	}

	wall, err := calendarWallTime(ctx, date, y, time.Month(monthIndex+1), d)
	if err != nil {
		return time.Time{}, err
	}

	// Follow Date's zone selection, checking subtraction before it can wrap.
	// In a DST gap the selected offset may differ from the result's Zone().
	zone := wall.In(date.Location())
	_, offset := zone.Zone()
	start, end := zone.ZoneBounds()
	seconds, err := calendarUnixSeconds(ctx, wall.Unix(), offset)
	if err != nil {
		return time.Time{}, err
	}

	provisional := time.Unix(seconds, 0)
	if (!start.IsZero() && provisional.Before(start)) || (!end.IsZero() && !provisional.Before(end)) {
		_, offset = provisional.In(date.Location()).Zone()
		seconds, err = calendarUnixSeconds(ctx, wall.Unix(), offset)
		if err != nil {
			return time.Time{}, err
		}
	}

	result := date.AddDate(nativeYears, nativeMonths, nativeDays)
	if result.Unix() != seconds || result.Nanosecond() != date.Nanosecond() {
		return time.Time{}, calendarRangeError()
	}

	return result, nil
}

func calendarWallTime(ctx context.Context, date time.Time, year int, month time.Month, day int) (time.Time, error) {
	h, min, sec := date.Clock()
	// This interval and its day normalization fit both native int widths and
	// time.Time comfortably. Keep ordinary calendar range checks allocation-free.
	if year >= -10000 && year <= 10000 && day >= -3660000 && day <= 3660000 {
		return time.Date(year, month, day, h, min, sec, date.Nanosecond(), time.UTC), nil
	}

	// A Gregorian 400-year cycle has exactly 146097 days. Reduce both the year
	// and day before asking Go to normalize, so a wrapped time.Date result
	// cannot masquerade as a representable date.
	dayIndex := int64(day) - 1
	cycles := int64(year)/400 + dayIndex/146097
	reference := time.Date(year%400, month, int(dayIndex%146097)+1, h, min, sec, date.Nanosecond(), time.UTC)
	seconds := new(big.Int).Mul(big.NewInt(cycles), big.NewInt(146097*86400))
	seconds.Add(seconds, big.NewInt(reference.Unix()))
	if !seconds.IsInt64() {
		return time.Time{}, calendarRangeError()
	}

	// Once the seconds fit, the reconstructed year also fits int64.
	normalizedYear := cycles*400 + int64(reference.Year())
	if _, ok := runtime.ToNativeInt(runtime.Int(normalizedYear)); !ok {
		return time.Time{}, calendarRangeError()
	}

	wall, err := runtime.ToDateTimeEpoch(ctx, runtime.NewInt64(seconds.Int64()), runtime.String("s"))
	if err != nil {
		return time.Time{}, calendarRangeError()
	}

	return wall.Time, nil
}

func calendarUnixSeconds(ctx context.Context, wall int64, offset int) (int64, error) {
	zoneOffset := int64(offset)
	if (zoneOffset > 0 && wall < math.MinInt64+zoneOffset) || (zoneOffset < 0 && wall > math.MaxInt64+zoneOffset) {
		return 0, calendarRangeError()
	}

	seconds := wall - zoneOffset
	// This inner interval is representable even with time.Time's epoch offset.
	// Extreme epochs use the runtime's authoritative DateTime range validation.
	if seconds < -(1<<50) || seconds > 1<<50 {
		if _, err := runtime.ToDateTimeEpoch(ctx, runtime.NewInt64(seconds), runtime.String("s")); err != nil {
			return 0, calendarRangeError()
		}
	}

	return seconds, nil
}

// Check the native component additions used by AddDate and month normalization.
func calendarComponent(base, delta int) (int, bool) {
	if (delta > 0 && base > math.MaxInt-delta) || (delta < 0 && base < math.MinInt-delta) {
		return 0, false
	}

	return base + delta, true
}

func calendarRangeError() error {
	return runtime.Error(runtime.ErrRange, "calendar shift exceeds the supported DateTime or native integer range")
}
