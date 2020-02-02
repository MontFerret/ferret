package datetime

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Unit specifies an unit of time (Millisecond, Second...).
type Unit int

const (
	Millisecond Unit = iota
	Second
	Minute
	Hour
	Day
	Week
	Month
	Year
)

var nanoseconds = []float64{
	1e6,
	1e9,
	6e10,
	36e11,
	864e11,
	6048e11,
	26784e11,
	31536e12,
}

// Nanosecond returns representation of an Unit
// in nanosconds
func (u Unit) Nanosecond() float64 {
	return nanoseconds[u]
}

// IsDatesEqual check if two partial dates match.
// This case the day means not the amount of days in Time,
// but the day of the month.
// The same rules applied to each unit.
func IsDatesEqual(tm1, tm2 time.Time, u Unit) bool {
	switch u {
	case Millisecond:
		tm1Msec := tm1.Nanosecond() / 1e6
		tm2Msec := tm2.Nanosecond() / 1e6
		return tm1Msec == tm2Msec
	case Second:
		return tm1.Second() == tm2.Second()
	case Minute:
		return tm1.Minute() == tm2.Minute()
	case Hour:
		return tm1.Hour() == tm2.Hour()
	case Day:
		return tm1.Day() == tm2.Day()
	case Week:
		tm1Wk := tm1.Day() / 7
		tm2Wk := tm2.Day() / 7
		return tm1Wk == tm2Wk
	case Month:
		return tm1.Month() == tm2.Month()
	case Year:
		return tm1.Year() == tm2.Year()
	}
	return false
}

// AddUnit add amount given in u to tm
func AddUnit(tm time.Time, amount int, u Unit) (res time.Time) {
	if u < Day {
		return tm.Add(time.Duration(amount) * time.Duration(int64(u.Nanosecond())))
	}

	switch u {
	case Day:
		res = tm.AddDate(0, 0, amount*1)
	case Week:
		res = tm.AddDate(0, 0, amount*7)
	case Month:
		res = tm.AddDate(0, amount*1, 0)
	case Year:
		res = tm.AddDate(amount*1, 0, 0)
	}

	return
}

// UnitFromString returns true and an Unit object if
// Unit with that name exists. Returns false, otherwise.
func UnitFromString(s string) (Unit, error) {
	switch strings.ToLower(s) {
	case "y", "year", "years":
		return Year, nil
	case "m", "month", "months":
		return Month, nil
	case "w", "week", "weeks":
		return Week, nil
	case "d", "day", "days":
		return Day, nil
	case "h", "hour", "hours":
		return Hour, nil
	case "i", "minute", "minutes":
		return Minute, nil
	case "s", "second", "seconds":
		return Second, nil
	case "f", "millisecond", "milliseconds":
		return Millisecond, nil
	}
	return -1, errors.Errorf("no such unit '%s'", s)
}
