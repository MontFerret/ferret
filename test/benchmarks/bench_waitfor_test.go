package benchmarks_test

import "testing"

const waitForValuePresentQuery = `
RETURN WAITFOR VALUE @candidate`

func BenchmarkWaitForValuePresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForValuePresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}
