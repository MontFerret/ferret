package compiler_test

import (
	"testing"
)

func BenchmarkStringLiteral(b *testing.B) {
	RunBenchmark(b, `
			RETURN "
FOO
BAR
"
		`)
}

func BenchmarkEmptyArray(b *testing.B) {
	RunBenchmark(b, `RETURN []`)
}

func BenchmarkStaticArray(b *testing.B) {
	RunBenchmark(b, `RETURN [1,2,3,4,5,6,7,8,9,10]`)
}

func BenchmarkEmptyObject(b *testing.B) {
	RunBenchmark(b, `RETURN {}`)
}

func BenchmarkUnaryOperatorExcl(b *testing.B) {
	RunBenchmark(b, `RETURN !TRUE`)
}

func BenchmarkUnaryOperatorQ(b *testing.B) {
	RunBenchmark(b, `
			LET foo = TRUE
			RETURN !foo ? TRUE : FALSE
		`)
}

func BenchmarkUnaryOperatorN(b *testing.B) {
	RunBenchmark(b, `
			LET v = 1
			RETURN -v
		`)
}

func BenchmarkTernaryOperator(b *testing.B) {
	RunBenchmark(b, `
			LET a = "a"
			LET b = "b"
			LET c = FALSE
			RETURN c ? a : b;
				
		`)
}

func BenchmarkTernaryOperatorDef(b *testing.B) {
	RunBenchmark(b, `
			LET a = "a"
			LET b = "b"
			LET c = FALSE
			RETURN c ? : a;
				
		`)
}

func BenchmarkForEmpty(b *testing.B) {
	RunBenchmark(b, `
			FOR i IN []
				RETURN i
		`)
}

func BenchmarkForStaticArray(b *testing.B) {
	RunBenchmark(b, `
			FOR i IN [1,2,3,4,5,6,7,8,9,10]
				RETURN i
		`)
}

func BenchmarkForRange(b *testing.B) {
	RunBenchmark(b, `
			FOR i IN 1..10
				RETURN i
		`)
}

func BenchmarkForObject(b *testing.B) {
	RunBenchmark(b, `
			FOR i IN {"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9":9, "10":10}
				RETURN i
		`)
}

func BenchmarkForNested(b *testing.B) {
	RunBenchmark(b, `
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					RETURN {[prop]: val}
		`)
}

func BenchmarkForTernary(b *testing.B) {
	RunBenchmark(b, `
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`)
}
