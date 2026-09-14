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
		Array(`RETURN [math::clamp(5,0,10),math::clamp(-1,0,10),math::clamp(20,0,10),math::clamp(2,3,3),math::clamp(0.25,0,1)]`, []any{5, 0, 10, 3, 0.25}),
		Array(`RETURN [math::sign(-12),math::sign(0),math::sign(8.5),math::trunc(3.9),math::trunc(-3.9),math::cbrt(27),math::cbrt(-8),math::hypot(3,4)]`, []any{-1, 0, 1, 3, -3, 3, -2, 5}),
		S(`RETURN math::e()==math::exp(1) AND math::log1p(0)==0 AND math::expm1(0)==0`, true),
		Nil(`LET x=0.0000000000000001 RETURN t::approx(math::log1p(x)/x,1,0.000001)`),
		Nil(`LET x=-0.0000000000000001 RETURN t::approx(math::log1p(x)/x,1,0.000001)`),
		Nil(`LET x=0.0000000000000001 RETURN t::approx(math::expm1(x)/x,1,0.000001)`),
		Nil(`LET x=-0.0000000000000001 RETURN t::approx(math::expm1(x)/x,1,0.000001)`),
		S(`RETURN IS_INT(math::clamp(5,0,10)) AND IS_INT(math::sign(1.5)) AND IS_INT(math::trunc(9007199254740993)) AND math::trunc(9007199254740993)==9007199254740993 AND IS_FLOAT(math::trunc(3.9))`, true),
		S(`RETURN IS_INT(math::clamp(5,0.5,10.5)) AND IS_FLOAT(math::clamp(5.5,0,10)) AND IS_INT(math::clamp(-0.5,0,10.5)) AND IS_FLOAT(math::clamp(-1,0.5,10)) AND IS_INT(math::clamp(20.5,0.5,10)) AND IS_FLOAT(math::clamp(20,0,10.5))`, true),
		S(`RETURN IS_INT(math::clamp(1,1.0,10.0)) AND IS_FLOAT(math::clamp(10.0,0,10)) AND IS_INT(math::clamp(3,3.0,3)) AND IS_FLOAT(math::clamp(3.0,3,3.0))`, true),
		S(`RETURN math::clamp(9007199254740993,0,9223372036854775807)==9007199254740993 AND math::clamp(9007199254740992.0,9007199254740993,9223372036854775807)==9007199254740993 AND math::clamp(9007199254740994.0,0,9007199254740993)==9007199254740993`, true),
		S(`LET nan=math::sqrt(-1) LET positive=math::exp(1000) LET negative=math::log(0) RETURN IS_NAN(math::clamp(nan,positive,positive)) AND IS_NAN(math::clamp(nan,negative,negative))`, true),
		S(`RETURN IS_NAN(math::log1p(-2)) AND IS_NAN(math::trunc(math::sqrt(-1))) AND IS_NAN(math::clamp(math::sqrt(-1),0,1))`, true),
		Array(`RETURN [MATH::CLAMP(5,0,10),MATH::SIGN(-1),MATH::TRUNC(-1.5),MATH::CBRT(-8),MATH::HYPOT(3,4),MATH::LOG1P(0),MATH::EXPM1(0),MATH::E()==math::e()]`, []any{5, -1, -1, -2, 5, 0, 0, true}),
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

	for _, test := range []struct {
		query    string
		category string
		position int
	}{
		{`RETURN math::clamp(5,10,0)`, "invalid argument", 2},
		{`RETURN math::clamp(5,9007199254740993,9007199254740992)`, "invalid argument", 2},
		{`RETURN math::clamp(5,math::sqrt(-1),10)`, "invalid argument", 2},
		{`RETURN math::clamp(5,0,math::sqrt(-1))`, "invalid argument", 3},
		{`RETURN math::sign(math::sqrt(-1))`, "invalid argument", 1},
		{`RETURN math::clamp("5",0,10)`, "invalid type", 1},
		{`RETURN math::clamp(5,"0",10)`, "invalid type", 2},
		{`RETURN math::clamp(5,0,"10")`, "invalid type", 3},
		{`RETURN math::hypot("3",4)`, "invalid type", 1},
		{`RETURN math::hypot(3,"4")`, "invalid type", 2},
	} {
		cases = append(cases, spec.NewSpec(test.query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Contains: []string{test.category, fmt.Sprintf("argument %d", test.position)}}))
	}

	for _, name := range []string{"sign", "trunc", "cbrt", "log1p", "expm1"} {
		cases = append(cases, spec.NewSpec(fmt.Sprintf(`RETURN math::%s("1")`, name)).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Contains: []string{"invalid type", "argument 1"}}))
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
