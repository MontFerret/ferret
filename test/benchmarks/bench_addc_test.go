package benchmarks_test

import (
	"testing"
	"time"
)

const (
	addConstNumericQuery = `
LET base = 1
RETURN FOR i IN 1..1000
  RETURN base + 2
`

	addConstNumericWithParamQuery = `
RETURN FOR i IN 1..1000
  RETURN @base + 2
`

	addConstStringQuery = `
LET foo = "bar"
RETURN FOR i IN 1..1000
  RETURN foo + " baz"`

	addConstStringWithParamQuery = `
RETURN FOR i IN 1..1000
  RETURN @foo + " baz"
`

	templateLiteralSimpleQuery = "RETURN FOR i IN 1..1000 RETURN `Hello ${@name}!`"

	templateLiteralNumericQuery = "RETURN FOR i IN 1..1000 RETURN `sum=${@a + @b}`"

	durationAddQuery = `
RETURN FOR i IN 1..1000
  RETURN @base + 2ms
`

	durationLiteralQuery = `RETURN FOR i IN 1..1000 RETURN 1.5s`

	durationExplicitAddQuery = `
RETURN FOR i IN 1..1000
  RETURN @base + TO_DURATION("2ms")
`

	durationExplicitCompareQuery = `
RETURN FOR i IN 1..1000
  RETURN @base > TO_DURATION(500)
`

	durationStrictCompareQuery = `
RETURN FOR i IN 1..1000
  RETURN @base > 500ms
`

	numericEqualityQuery = `
RETURN FOR i IN 1..1000
  RETURN @base == 1
`

	durationExplicitEqualityQuery = `
RETURN FOR i IN 1..1000
  RETURN @base == TO_DURATION("1s")
`

	durationStrictEqualityQuery = `
RETURN FOR i IN 1..1000
  RETURN @base == 1s
`

	equalityJumpConstQuery = `
RETURN FOR i IN 1..1000
  RETURN @left == 1 ? i : 0
`

	equalityJumpRegisterQuery = `
RETURN FOR i IN 1..1000
  RETURN @left == @right ? i : 0
`

	quantifiedComparisonQuery = `
RETURN FOR i IN 1..1000
  RETURN @values ANY > @threshold
`

	dateTimeAddQuery = `
RETURN FOR i IN 1..1000
  RETURN @base + 2ms
`

	dateTimeConversionQuery = `
RETURN FOR i IN 1..1000
  RETURN TO_DATETIME(@value)
`

	dateTimeEpochConversionQuery = `
RETURN FOR i IN 1..1000
  RETURN TO_DATETIME(@value, "s")
`
)

func BenchmarkAddNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstNumericQuery)
}

func BenchmarkAddNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstNumericQuery)
}

func BenchmarkAddConstNumericWithParam_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstNumericWithParam_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstString_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstStringQuery)
}

func BenchmarkAddConstString_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstStringQuery)
}

func BenchmarkAddConstStringWithParam_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkAddConstStringWithParam_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkTemplateLiteralSimple_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralSimple_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkTemplateLiteralNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkDurationAdd_O0(b *testing.B) {
	RunBenchmarkO0(b, durationAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationAdd_O1(b *testing.B) {
	RunBenchmarkO1(b, durationAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationLiteral_O0(b *testing.B) {
	RunBenchmarkO0(b, durationLiteralQuery)
}

func BenchmarkDurationLiteral_O1(b *testing.B) {
	RunBenchmarkO1(b, durationLiteralQuery)
}

func BenchmarkDurationExplicitAdd_O0(b *testing.B) {
	RunBenchmarkO0(b, durationExplicitAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitAdd_O1(b *testing.B) {
	RunBenchmarkO1(b, durationExplicitAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitCompare_O0(b *testing.B) {
	RunBenchmarkO0(b, durationExplicitCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitCompare_O1(b *testing.B) {
	RunBenchmarkO1(b, durationExplicitCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictCompare_O0(b *testing.B) {
	RunBenchmarkO0(b, durationStrictCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictCompare_O1(b *testing.B) {
	RunBenchmarkO1(b, durationStrictCompareQuery, WithParam("base", time.Second))
}

func BenchmarkNumericEquality_O0(b *testing.B) {
	RunBenchmarkO0(b, numericEqualityQuery, WithParam("base", 1))
}

func BenchmarkNumericEquality_O1(b *testing.B) {
	RunBenchmarkO1(b, numericEqualityQuery, WithParam("base", 1))
}

func BenchmarkDurationExplicitEquality_O0(b *testing.B) {
	RunBenchmarkO0(b, durationExplicitEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitEquality_O1(b *testing.B) {
	RunBenchmarkO1(b, durationExplicitEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictEquality_O0(b *testing.B) {
	RunBenchmarkO0(b, durationStrictEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictEquality_O1(b *testing.B) {
	RunBenchmarkO1(b, durationStrictEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkEqualityJumpConst_O0(b *testing.B) {
	RunBenchmarkO0(b, equalityJumpConstQuery, WithParam("left", 1))
}

func BenchmarkEqualityJumpConst_O1(b *testing.B) {
	RunBenchmarkO1(b, equalityJumpConstQuery, WithParam("left", 1))
}

func BenchmarkEqualityJumpRegister_O0(b *testing.B) {
	RunBenchmarkO0(b, equalityJumpRegisterQuery, WithParam("left", 1), WithParam("right", 1))
}

func BenchmarkEqualityJumpRegister_O1(b *testing.B) {
	RunBenchmarkO1(b, equalityJumpRegisterQuery, WithParam("left", 1), WithParam("right", 1))
}

func BenchmarkQuantifiedComparison_O0(b *testing.B) {
	RunBenchmarkO0(b, quantifiedComparisonQuery, WithParam("values", []any{1, 2, 3, 4, 5, 6, 7, 8}), WithParam("threshold", 7))
}

func BenchmarkQuantifiedComparison_O1(b *testing.B) {
	RunBenchmarkO1(b, quantifiedComparisonQuery, WithParam("values", []any{1, 2, 3, 4, 5, 6, 7, 8}), WithParam("threshold", 7))
}

func BenchmarkDateTimeAdd_O0(b *testing.B) {
	RunBenchmarkO0(b, dateTimeAddQuery, WithParam("base", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)))
}

func BenchmarkDateTimeAdd_O1(b *testing.B) {
	RunBenchmarkO1(b, dateTimeAddQuery, WithParam("base", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)))
}

func BenchmarkDateTimeConversion_O0(b *testing.B) {
	RunBenchmarkO0(b, dateTimeConversionQuery, WithParam("value", "2026-08-02T12:00:00Z"))
}

func BenchmarkDateTimeConversion_O1(b *testing.B) {
	RunBenchmarkO1(b, dateTimeConversionQuery, WithParam("value", "2026-08-02T12:00:00Z"))
}

func BenchmarkDateTimeEpochConversion_O0(b *testing.B) {
	RunBenchmarkO0(b, dateTimeEpochConversionQuery, WithParam("value", 1_690_992_000))
}

func BenchmarkDateTimeEpochConversion_O1(b *testing.B) {
	RunBenchmarkO1(b, dateTimeEpochConversionQuery, WithParam("value", 1_690_992_000))
}
