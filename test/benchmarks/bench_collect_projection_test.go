package benchmarks_test

import "testing"

const (
	collectProjectionAllVars = `
FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups
  RETURN groups
`
	collectProjectionSingleGroup = `
FOR i IN 1..10000
  COLLECT g = "only" INTO groups
  RETURN groups
`

	collectProjectionKeep = `
FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups KEEP a, b
  RETURN groups
`

	collectProjectionCustom = `
FOR i IN 1..200
  LET a = i + 1
  LET b = i + 2
  COLLECT g = i % 5 INTO groups = { a: a, b: b }
  RETURN groups
`

	collectProjectionCount = `
FOR i IN 1..200
  COLLECT WITH COUNT INTO total
  RETURN total
`
)

func BenchmarkCollectProjection_AllVars_O0(b *testing.B) {
	RunBenchmarkO0(b, collectProjectionAllVars)
}

func BenchmarkCollectProjection_AllVars_O1(b *testing.B) {
	RunBenchmarkO1(b, collectProjectionAllVars)
}

func BenchmarkCollectProjection_SingleGroup_O0(b *testing.B) {
	RunBenchmarkO0(b, collectProjectionSingleGroup)
}

func BenchmarkCollectProjection_SingleGroup_O1(b *testing.B) {
	RunBenchmarkO1(b, collectProjectionSingleGroup)
}

func BenchmarkCollectProjection_Keep_O0(b *testing.B) {
	RunBenchmarkO0(b, collectProjectionKeep)
}

func BenchmarkCollectProjection_Keep_O1(b *testing.B) {
	RunBenchmarkO1(b, collectProjectionKeep)
}

func BenchmarkCollectProjection_Custom_O0(b *testing.B) {
	RunBenchmarkO0(b, collectProjectionCustom)
}

func BenchmarkCollectProjection_Custom_O1(b *testing.B) {
	RunBenchmarkO1(b, collectProjectionCustom)
}

func BenchmarkCollectProjection_Count_O0(b *testing.B) {
	RunBenchmarkO0(b, collectProjectionCount)
}

func BenchmarkCollectProjection_Count_O1(b *testing.B) {
	RunBenchmarkO1(b, collectProjectionCount)
}
