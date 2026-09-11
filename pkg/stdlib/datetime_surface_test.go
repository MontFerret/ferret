package stdlib_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestDateTimeSurface(t *testing.T) {
	want := map[string][]int{
		"date":                    {1, 2},
		"date_add":                {3},
		"date_compare":            {3, 4},
		"date_day":                {1},
		"date_dayofweek":          {1},
		"date_dayofyear":          {1},
		"date_days_in_month":      {1},
		"date_diff":               {3, 4},
		"date_format":             {2},
		"date_hour":               {1},
		"date_leapyear":           {1},
		"date_millisecond":        {1},
		"date_minute":             {1},
		"date_month":              {1},
		"date_quarter":            {1},
		"date_second":             {1},
		"date_subtract":           {3},
		"date_year":               {1},
		"datetime::add":           {3},
		"datetime::day":           {1},
		"datetime::day_of_week":   {1},
		"datetime::day_of_year":   {1},
		"datetime::days_in_month": {1},
		"datetime::diff":          {3},
		"datetime::format":        {2},
		"datetime::hour":          {1},
		"datetime::is_leap_year":  {1},
		"datetime::millisecond":   {1},
		"datetime::minute":        {1},
		"datetime::month":         {1},
		"datetime::now":           {0},
		"datetime::parse":         {1, 2},
		"datetime::quarter":       {1},
		"datetime::same":          {3},
		"datetime::second":        {1},
		"datetime::subtract":      {3},
		"datetime::year":          {1},
		"now":                     {0},
	}
	names := make([]string, 0, len(want))
	for name := range want {
		names = append(names, name)
	}

	slices.Sort(names)
	got := functionNames(buildFunctions(t, stdlib.Only(stdlib.DateTime)))
	if !reflect.DeepEqual(got, names) {
		t.Fatalf("datetime surface = %v, want %v", got, names)
	}

	for _, set := range []stdlib.Set{stdlib.Only(stdlib.DateTime), stdlib.Full(), stdlib.Safe()} {
		functions := buildFunctions(t, set)
		without := buildFunctions(t, set.Without(stdlib.DateTime))
		for name, arities := range want {
			for _, spelling := range []string{name, strings.ToUpper(name)} {
				for arity := 0; arity <= 4; arity++ {
					if hasFixedArity(functions, spelling, arity) != slices.Contains(arities, arity) || functions.Var().Has(spelling) {
						t.Fatalf("wrong arity %d for %s", arity, spelling)
					}
				}

				if without.Has(spelling) {
					t.Fatalf("datetime group removal retained %s", spelling)
				}
			}
		}

		for _, name := range []string{"date_now", "datetime::compare", "datetime::date", "datetime::date_year", "datetime::unix", "datetime::start_of"} {
			if functions.Has(name) {
				t.Fatalf("unexpected function %s", name)
			}
		}
	}
}
