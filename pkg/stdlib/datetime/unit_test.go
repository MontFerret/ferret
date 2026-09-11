package datetime

import (
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestUnitParsing(t *testing.T) {
	for expected, names := range map[unit][]string{
		millisecond: {"millisecond", "milliseconds", "f"},
		second:      {"second", "seconds", "s"},
		minute:      {"minute", "minutes", "i"},
		hour:        {"hour", "hours", "h"},
		day:         {"day", "days", "d"},
		week:        {"week", "weeks", "w"},
		month:       {"month", "months", "m"},
		year:        {"year", "years", "y"},
	} {
		for _, name := range names {
			for _, spelling := range []string{name, strings.ToUpper(name)} {
				got, err := parseUnit(runtime.String(spelling), 2)
				if err != nil || got != expected {
					t.Fatalf("%s = %v, %v; want %v", spelling, got, err, expected)
				}
			}
		}
	}

	for _, input := range []runtime.Value{runtime.String(""), runtime.String("unknown"), runtime.String(" hour "), runtime.None, runtime.Int(1)} {
		_, err := parseUnit(input, 3)
		pos, ok, _ := runtime.InvalidArgumentDetails(err)
		if !errors.Is(err, runtime.ErrInvalidArgument) || !ok || pos != 3 {
			t.Fatalf("%v: expected argument 3 error, got %v", input, err)
		}
	}
}
