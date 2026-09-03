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

func BenchmarkAddNumeric_None(b *testing.B) {
	RunBenchmarkNone(b, addConstNumericQuery)
}

func BenchmarkAddNumeric_Basic(b *testing.B) {
	RunBenchmarkBasic(b, addConstNumericQuery)
}

func BenchmarkAddNumeric_Full(b *testing.B) {
	RunBenchmarkFull(b, addConstNumericQuery)
}

func BenchmarkAddConstNumericWithParam_None(b *testing.B) {
	RunBenchmarkNone(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstNumericWithParam_Basic(b *testing.B) {
	RunBenchmarkBasic(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstNumericWithParam_Full(b *testing.B) {
	RunBenchmarkFull(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstString_None(b *testing.B) {
	RunBenchmarkNone(b, addConstStringQuery)
}

func BenchmarkAddConstString_Basic(b *testing.B) {
	RunBenchmarkBasic(b, addConstStringQuery)
}

func BenchmarkAddConstString_Full(b *testing.B) {
	RunBenchmarkFull(b, addConstStringQuery)
}

func BenchmarkAddConstStringWithParam_None(b *testing.B) {
	RunBenchmarkNone(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkAddConstStringWithParam_Basic(b *testing.B) {
	RunBenchmarkBasic(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkAddConstStringWithParam_Full(b *testing.B) {
	RunBenchmarkFull(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkTemplateLiteralSimple_None(b *testing.B) {
	RunBenchmarkNone(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralSimple_Basic(b *testing.B) {
	RunBenchmarkBasic(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralSimple_Full(b *testing.B) {
	RunBenchmarkFull(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralNumeric_None(b *testing.B) {
	RunBenchmarkNone(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkTemplateLiteralNumeric_Basic(b *testing.B) {
	RunBenchmarkBasic(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkTemplateLiteralNumeric_Full(b *testing.B) {
	RunBenchmarkFull(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkDurationAdd_None(b *testing.B) {
	RunBenchmarkNone(b, durationAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationAdd_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationAdd_Full(b *testing.B) {
	RunBenchmarkFull(b, durationAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationLiteral_None(b *testing.B) {
	RunBenchmarkNone(b, durationLiteralQuery)
}

func BenchmarkDurationLiteral_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationLiteralQuery)
}

func BenchmarkDurationLiteral_Full(b *testing.B) {
	RunBenchmarkFull(b, durationLiteralQuery)
}

func BenchmarkDurationExplicitAdd_None(b *testing.B) {
	RunBenchmarkNone(b, durationExplicitAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitAdd_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationExplicitAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitAdd_Full(b *testing.B) {
	RunBenchmarkFull(b, durationExplicitAddQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitCompare_None(b *testing.B) {
	RunBenchmarkNone(b, durationExplicitCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitCompare_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationExplicitCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitCompare_Full(b *testing.B) {
	RunBenchmarkFull(b, durationExplicitCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictCompare_None(b *testing.B) {
	RunBenchmarkNone(b, durationStrictCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictCompare_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationStrictCompareQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictCompare_Full(b *testing.B) {
	RunBenchmarkFull(b, durationStrictCompareQuery, WithParam("base", time.Second))
}

func BenchmarkNumericEquality_None(b *testing.B) {
	RunBenchmarkNone(b, numericEqualityQuery, WithParam("base", 1))
}

func BenchmarkNumericEquality_Basic(b *testing.B) {
	RunBenchmarkBasic(b, numericEqualityQuery, WithParam("base", 1))
}

func BenchmarkNumericEquality_Full(b *testing.B) {
	RunBenchmarkFull(b, numericEqualityQuery, WithParam("base", 1))
}

func BenchmarkDurationExplicitEquality_None(b *testing.B) {
	RunBenchmarkNone(b, durationExplicitEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitEquality_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationExplicitEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationExplicitEquality_Full(b *testing.B) {
	RunBenchmarkFull(b, durationExplicitEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictEquality_None(b *testing.B) {
	RunBenchmarkNone(b, durationStrictEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictEquality_Basic(b *testing.B) {
	RunBenchmarkBasic(b, durationStrictEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkDurationStrictEquality_Full(b *testing.B) {
	RunBenchmarkFull(b, durationStrictEqualityQuery, WithParam("base", time.Second))
}

func BenchmarkEqualityJumpConst_None(b *testing.B) {
	RunBenchmarkNone(b, equalityJumpConstQuery, WithParam("left", 1))
}

func BenchmarkEqualityJumpConst_Basic(b *testing.B) {
	RunBenchmarkBasic(b, equalityJumpConstQuery, WithParam("left", 1))
}

func BenchmarkEqualityJumpConst_Full(b *testing.B) {
	RunBenchmarkFull(b, equalityJumpConstQuery, WithParam("left", 1))
}

func BenchmarkEqualityJumpRegister_None(b *testing.B) {
	RunBenchmarkNone(b, equalityJumpRegisterQuery, WithParam("left", 1), WithParam("right", 1))
}

func BenchmarkEqualityJumpRegister_Basic(b *testing.B) {
	RunBenchmarkBasic(b, equalityJumpRegisterQuery, WithParam("left", 1), WithParam("right", 1))
}

func BenchmarkEqualityJumpRegister_Full(b *testing.B) {
	RunBenchmarkFull(b, equalityJumpRegisterQuery, WithParam("left", 1), WithParam("right", 1))
}

func BenchmarkQuantifiedComparison_None(b *testing.B) {
	RunBenchmarkNone(b, quantifiedComparisonQuery, WithParam("values", []any{1, 2, 3, 4, 5, 6, 7, 8}), WithParam("threshold", 7))
}

func BenchmarkQuantifiedComparison_Basic(b *testing.B) {
	RunBenchmarkBasic(b, quantifiedComparisonQuery, WithParam("values", []any{1, 2, 3, 4, 5, 6, 7, 8}), WithParam("threshold", 7))
}

func BenchmarkQuantifiedComparison_Full(b *testing.B) {
	RunBenchmarkFull(b, quantifiedComparisonQuery, WithParam("values", []any{1, 2, 3, 4, 5, 6, 7, 8}), WithParam("threshold", 7))
}

func BenchmarkDateTimeAdd_None(b *testing.B) {
	RunBenchmarkNone(b, dateTimeAddQuery, WithParam("base", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)))
}

func BenchmarkDateTimeAdd_Basic(b *testing.B) {
	RunBenchmarkBasic(b, dateTimeAddQuery, WithParam("base", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)))
}

func BenchmarkDateTimeAdd_Full(b *testing.B) {
	RunBenchmarkFull(b, dateTimeAddQuery, WithParam("base", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)))
}

func BenchmarkDateTimeConversion_None(b *testing.B) {
	RunBenchmarkNone(b, dateTimeConversionQuery, WithParam("value", "2026-08-02T12:00:00Z"))
}

func BenchmarkDateTimeConversion_Basic(b *testing.B) {
	RunBenchmarkBasic(b, dateTimeConversionQuery, WithParam("value", "2026-08-02T12:00:00Z"))
}

func BenchmarkDateTimeConversion_Full(b *testing.B) {
	RunBenchmarkFull(b, dateTimeConversionQuery, WithParam("value", "2026-08-02T12:00:00Z"))
}

func BenchmarkDateTimeEpochConversion_None(b *testing.B) {
	RunBenchmarkNone(b, dateTimeEpochConversionQuery, WithParam("value", 1_690_992_000))
}

func BenchmarkDateTimeEpochConversion_Basic(b *testing.B) {
	RunBenchmarkBasic(b, dateTimeEpochConversionQuery, WithParam("value", 1_690_992_000))
}

func BenchmarkDateTimeEpochConversion_Full(b *testing.B) {
	RunBenchmarkFull(b, dateTimeEpochConversionQuery, WithParam("value", 1_690_992_000))
}
