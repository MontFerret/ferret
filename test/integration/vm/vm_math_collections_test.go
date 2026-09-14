package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMathCollectionFunctions(t *testing.T) {
	cases := []spec.Spec{
		Array(`RETURN [MEDIAN([1,2,3]), MEDIAN([1,2,3,4]), MEDIAN([1,2,3,100])]`, []any{2, 2.5, 2.5}),
		Array(`RETURN [SUM([-3,1.5,2]), AVERAGE([-3,1.5,2]), MIN([-3,1.5,2]), MAX([-3,1.5,2])]`, []any{0.5, 0.5 / 3, -3, 2}),
		Array(`RETURN [VARIANCE_POPULATION([1,3]), VARIANCE_SAMPLE([1,3]), STDDEV_POPULATION([1,3]), STDDEV_SAMPLE([1,3])]`, []any{1, 2, 1, 1.4142135623730951}),
		Array(`RETURN [SUM([]), AVERAGE([]), MIN([]), MAX([]), MEDIAN([])]`, []any{0, 0, nil, nil, nil}),
		Array(`RETURN [PERCENTILE([100,3,1,2],1), PERCENTILE([100,3,1,2],100), PERCENTILE([100,3,1,2],50,"interpolation"), PERCENTILE([100,3,1,2],50,"unknown")]`, []any{1, 100, 2.5, 2}),
		Array(`LET values = [100,3,1,2] LET median = MEDIAN(values) LET percentile = PERCENTILE(values,50) RETURN [median,percentile,values]`, []any{2.5, 2, []any{100, 3, 1, 2}}),
		S(`RETURN IS_NAN(PERCENTILE([], "invalid"))`, true),
		S(`RETURN IS_NAN(VARIANCE_SAMPLE([1])) AND IS_NAN(STDDEV_SAMPLE([1]))`, true),
		Array(`RETURN [VARIANCE_POPULATION([1]), STDDEV_POPULATION([1])]`, []any{0, 0}),
	}
	for _, test := range []struct {
		want      any
		empty     any
		name      string
		extra     string
		undefined bool
	}{
		{name: "SUM", want: 4, empty: 0},
		{name: "AVERAGE", want: 2, empty: 0},
		{name: "MIN", want: 1},
		{name: "MAX", want: 3},
		{name: "MEDIAN", want: 2},
		{name: "VARIANCE_POPULATION", want: 1, undefined: true},
		{name: "VARIANCE_SAMPLE", want: 2, undefined: true},
		{name: "STDDEV_POPULATION", want: 1, undefined: true},
		{name: "STDDEV_SAMPLE", want: 1.4142135623730951, undefined: true},
		{name: "PERCENTILE", extra: ",50", want: 1, undefined: true},
		{name: "PERCENTILE", extra: `,50,"interpolation"`, want: 2, undefined: true},
	} {
		for _, skipped := range []string{`"2"`, "NONE", "TRUE"} {
			cases = append(cases, S(fmt.Sprintf("RETURN %s([1,%s,3]%s)", test.name, skipped, test.extra), test.want))
		}

		for _, input := range []string{`[]`, `["2",NONE,TRUE]`} {
			query := fmt.Sprintf("%s(%s%s)", test.name, input, test.extra)
			if test.undefined {
				cases = append(cases, S("RETURN IS_NAN("+query+")", true))
			} else {
				cases = append(cases, S("RETURN "+query, test.empty))
			}
		}
	}

	RunSpecs(t, cases)
}
