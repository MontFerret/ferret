package benchmarks_test

import "testing"

const (
	collectProjectionAllVars = `
RETURN FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups
  RETURN groups
`
	collectProjectionSingleGroup = `
RETURN FOR i IN 1..10000
  COLLECT g = "only" INTO groups
  RETURN groups
`

	collectProjectionKeep = `
RETURN FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups KEEP a, b
  RETURN groups
`

	collectProjectionCustom = `
RETURN FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups = { a: a, b: b }
  RETURN groups
`

	collectProjectionCount = `
RETURN FOR i IN 1..200
  COLLECT WITH COUNT INTO total
  RETURN total
`
)

func BenchmarkCollectProjection_AllVars_None(b *testing.B) {
	RunBenchmarkNone(b, collectProjectionAllVars)
}

func BenchmarkCollectProjection_AllVars_Basic(b *testing.B) {
	RunBenchmarkBasic(b, collectProjectionAllVars)
}

func BenchmarkCollectProjection_AllVars_Full(b *testing.B) {
	RunBenchmarkFull(b, collectProjectionAllVars)
}

func BenchmarkCollectProjection_SingleGroup_None(b *testing.B) {
	RunBenchmarkNone(b, collectProjectionSingleGroup)
}

func BenchmarkCollectProjection_SingleGroup_Basic(b *testing.B) {
	RunBenchmarkBasic(b, collectProjectionSingleGroup)
}

func BenchmarkCollectProjection_SingleGroup_Full(b *testing.B) {
	RunBenchmarkFull(b, collectProjectionSingleGroup)
}

func BenchmarkCollectProjection_Keep_None(b *testing.B) {
	RunBenchmarkNone(b, collectProjectionKeep)
}

func BenchmarkCollectProjection_Keep_Basic(b *testing.B) {
	RunBenchmarkBasic(b, collectProjectionKeep)
}

func BenchmarkCollectProjection_Keep_Full(b *testing.B) {
	RunBenchmarkFull(b, collectProjectionKeep)
}

func BenchmarkCollectProjection_Custom_None(b *testing.B) {
	RunBenchmarkNone(b, collectProjectionCustom)
}

func BenchmarkCollectProjection_Custom_Basic(b *testing.B) {
	RunBenchmarkBasic(b, collectProjectionCustom)
}

func BenchmarkCollectProjection_Custom_Full(b *testing.B) {
	RunBenchmarkFull(b, collectProjectionCustom)
}

func BenchmarkCollectProjection_Count_None(b *testing.B) {
	RunBenchmarkNone(b, collectProjectionCount)
}

func BenchmarkCollectProjection_Count_Basic(b *testing.B) {
	RunBenchmarkBasic(b, collectProjectionCount)
}

func BenchmarkCollectProjection_Count_Full(b *testing.B) {
	RunBenchmarkFull(b, collectProjectionCount)
}
