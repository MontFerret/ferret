package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestCanonicalMathFunctions(t *testing.T) {
	cases := []spec.Spec{
		Array(`RETURN [math::abs(-2),math::ceil(1.2),math::floor(1.9),math::round(1.9),math::sum([1,2,3]),math::mean([1,2,3]),math::min([1,2,3]),math::max([1,2,3]),math::median([1,2,3,100])]`, []any{2, 2, 1, 2, 6, 2, 1, 3, 2.5}),
		Array(`RETURN [math::variance([1,3]),math::variance_sample([1,3]),math::stddev([1,3]),math::stddev_sample([1,3])]`, []any{1, 2, 1, 1.4142135623730951}),
		Array(`RETURN [math::pow(2,3),math::sqrt(4),math::exp(0),math::exp2(3),math::log(1),math::log2(8),math::log10(100)]`, []any{8, 2, 1, 8, 0, 3, 2}),
		Array(`RETURN [math::sin(0),math::cos(0),math::tan(0),math::asin(0),math::acos(1),math::atan(0),math::atan2(0,1),math::degrees(math::pi()),math::radians(180)==math::pi()]`, []any{0, 1, 0, 0, 0, 0, 0, 180, true}),
		Array(`RETURN [math::sum([]),math::min([]),math::max([]),average([]),median([]),sum([1,"2",3])]`, []any{0, nil, nil, 0, nil, 4}),
		Array(`RETURN [arrays::range(1,4),range(1,4),arrays::range(4,1,-1),arrays::range(-0.5,0.5,0.5)]`, []any{[]any{1, 2, 3, 4}, []any{1, 2, 3, 4}, []any{4, 3, 2, 1}, []any{-0.5, 0, 0.5}}),
		Array(`RETURN [math::percentile([0,100,200,300,400],0),math::percentile([0,100,200,300,400],25),math::percentile([0,100,200,300,400],50),math::percentile([0,100,200,300,400],75),math::percentile([0,100,200,300,400],95),math::percentile([0,100,200,300,400],100)]`, []any{0, 100, 200, 300, 380, 400}),
		Nil(`RETURN t::approx(math::percentile([0,100,200,300,400],99.9),399.6,0.00001)`),
		Array(`LET v=[100,3,1,2] LET a=math::median(v) LET b=math::percentile(v,50) RETURN [a,b,v]`, []any{2.5, 2.5, []any{100, 3, 1, 2}}),
		Array(`RETURN [math::sum([1,"2",3])?,math::mean([1,none,3])?,math::percentile([],101)?]`, []any{nil, nil, nil}),
		Array(`RETURN FOR v IN [1,"2",3] COLLECT g=1 AGGREGATE a=math::sum(v)?,b=SUM(v) RETURN [a,b]`, []any{[]any{nil, 4}}),
		S(`RETURN IS_NAN(math::variance_sample([1])) AND IS_NAN(math::stddev_sample([1]))`, true),
		Array(`RETURN [percentile([1,2,3,100],50),percentile([1,2,3,100],50,"unknown"),percentile([1,2,3,100],50,"interpolation"),IS_NAN(percentile([],"invalid"))]`, []any{2, 2, 2.5, true}),
	}

	for _, name := range []string{"mean", "median", "variance", "variance_sample", "stddev", "stddev_sample", "percentile"} {
		extra := ""
		if name == "percentile" {
			extra = ",50"
		}

		cases = append(cases, S(fmt.Sprintf("RETURN IS_NAN(math::%s([]%s))", name, extra), true))
	}

	for _, name := range []string{"sum", "min", "max", "mean", "median", "variance", "variance_sample", "stddev", "stddev_sample", "percentile"} {
		extra := ""
		if name == "percentile" {
			extra = ",50"
		}

		cases = append(cases, spec.NewSpec(fmt.Sprintf(`RETURN math::%s([1,"2",3]%s)`, name, extra)).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Contains: []string{"invalid type", "index 1"}}))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, level.String(), c, cases)
	}
}
