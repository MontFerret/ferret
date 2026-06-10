package benchmarks_test

import "testing"

const bareArrayQuestionQuery = `RETURN @arr[?]`

const bareLocalArrayQuestionQuery = `
LET arr = [1, 2, 3, 4, 5, 6, 7, 8]
RETURN arr[?]`

var bareArrayQuestionValues = []any{1, 2, 3, 4, 5, 6, 7, 8}

func BenchmarkBareArrayQuestion_O0(b *testing.B) {
	RunBenchmarkO0(b, bareArrayQuestionQuery, WithParam("arr", bareArrayQuestionValues))
}

func BenchmarkBareArrayQuestion_O1(b *testing.B) {
	RunBenchmarkO1(b, bareArrayQuestionQuery, WithParam("arr", bareArrayQuestionValues))
}

func BenchmarkBareLocalArrayQuestion_O0(b *testing.B) {
	RunBenchmarkO0(b, bareLocalArrayQuestionQuery)
}

func BenchmarkBareLocalArrayQuestion_O1(b *testing.B) {
	RunBenchmarkO1(b, bareLocalArrayQuestionQuery)
}
