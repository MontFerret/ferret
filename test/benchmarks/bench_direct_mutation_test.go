package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const (
	compilerDirectMutationQuery = `
VAR count = 1
count += 1

LET obj = {
  count: count,
  profile: { nested: { value: 1 } },
  items: [{ count: 1 }, { count: 2 }]
}
obj.count = count
obj.nickname = "tim"

LET key = "theme"
obj[key] = "dark"

LET arr = [1, 2, 3]
arr[1] = obj.count

obj.items[0].count += 1
obj?.profile?.nested.value = arr[1]
obj?.missing?.count += 1
obj.transient = "drop"
DELETE obj.transient
DELETE obj?.missing?.legacy

RETURN obj.items[0].count + arr[1]
`

	directMutationBindingNumericQuery = `
VAR count = 1
count += 1
RETURN count
`

	directMutationObjectPropertyQuery = `
LET obj = { count: 1 }
obj.count = 2
obj.extra = 3
RETURN obj.count + obj.extra
`

	directMutationDynamicKeyQuery = `
LET obj = {}
LET key = "count"
obj[key] = 2
RETURN obj.count
`

	directMutationArrayIndexQuery = `
LET arr = [1, 2, 3]
arr[1] = 20
RETURN arr[1]
`

	directMutationNestedAugmentedQuery = `
LET data = { items: [{ count: 1 }] }
data.items[0].count += 1
RETURN data.items[0].count
`

	directMutationSafeNoopQuery = `
LET user = NONE
user?.profile?.count += 1
RETURN 1
`

	directMutationSafePresentQuery = `
LET user = { profile: { count: 1 } }
user?.profile?.count += 1
RETURN user.profile.count
`

	directMutationUDFCapturedRootQuery = `
LET user = { count: 1 }
FUNC inc() (
  user.count += 1
  RETURN user.count
)
RETURN inc()
`

	directMutationDeletePropertyQuery = `
LET obj = { count: 1, extra: 2 }
DELETE obj.extra
RETURN obj.count
`

	directMutationDeleteDynamicKeyQuery = `
LET obj = { debug: 1, keep: 2 }
LET key = "debug"
DELETE obj[key]
RETURN obj.keep
`

	directMutationDeleteSafeNoopQuery = `
LET obj = NONE
DELETE obj?.debug
RETURN 1
`
)

func BenchmarkCompilerCompileDirectMutation_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerDirectMutationQuery, compiler.O0)
}

func BenchmarkCompilerCompileDirectMutation_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerDirectMutationQuery, compiler.O1)
}

func BenchmarkDirectMutation_BindingNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationBindingNumericQuery)
}

func BenchmarkDirectMutation_BindingNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationBindingNumericQuery)
}

func BenchmarkDirectMutation_ObjectProperty_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationObjectPropertyQuery)
}

func BenchmarkDirectMutation_ObjectProperty_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationObjectPropertyQuery)
}

func BenchmarkDirectMutation_DynamicKey_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationDynamicKeyQuery)
}

func BenchmarkDirectMutation_DynamicKey_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationDynamicKeyQuery)
}

func BenchmarkDirectMutation_ArrayIndex_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationArrayIndexQuery)
}

func BenchmarkDirectMutation_ArrayIndex_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationArrayIndexQuery)
}

func BenchmarkDirectMutation_NestedAugmented_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationNestedAugmentedQuery)
}

func BenchmarkDirectMutation_NestedAugmented_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationNestedAugmentedQuery)
}

func BenchmarkDirectMutation_SafeNoop_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationSafeNoopQuery)
}

func BenchmarkDirectMutation_SafeNoop_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationSafeNoopQuery)
}

func BenchmarkDirectMutation_SafePresent_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationSafePresentQuery)
}

func BenchmarkDirectMutation_SafePresent_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationSafePresentQuery)
}

func BenchmarkDirectMutation_UDFCapturedRoot_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationUDFCapturedRootQuery)
}

func BenchmarkDirectMutation_UDFCapturedRoot_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationUDFCapturedRootQuery)
}

func BenchmarkDirectMutation_DeleteProperty_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationDeletePropertyQuery)
}

func BenchmarkDirectMutation_DeleteProperty_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationDeletePropertyQuery)
}

func BenchmarkDirectMutation_DeleteDynamicKey_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationDeleteDynamicKeyQuery)
}

func BenchmarkDirectMutation_DeleteDynamicKey_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationDeleteDynamicKeyQuery)
}

func BenchmarkDirectMutation_DeleteSafeNoop_O0(b *testing.B) {
	RunBenchmarkO0(b, directMutationDeleteSafeNoopQuery)
}

func BenchmarkDirectMutation_DeleteSafeNoop_O1(b *testing.B) {
	RunBenchmarkO1(b, directMutationDeleteSafeNoopQuery)
}
