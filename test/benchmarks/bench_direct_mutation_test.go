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
FUNC inc() {
  user.count += 1
  RETURN user.count
}
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

	directMutationDeleteArrayIndexQuery = `
LET arr = [1, 2, 3]
DELETE arr[1]
RETURN arr[1]
`

	directMutationDeleteSafeNoopQuery = `
LET obj = NONE
DELETE obj?.debug
RETURN 1
`
)

func BenchmarkCompilerCompileDirectMutation_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerDirectMutationQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileDirectMutation_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerDirectMutationQuery, compiler.OptimizationBasic)
}

func BenchmarkCompilerCompileDirectMutation_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerDirectMutationQuery, compiler.OptimizationFull)
}

func BenchmarkDirectMutation_BindingNumeric_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationBindingNumericQuery)
}

func BenchmarkDirectMutation_BindingNumeric_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationBindingNumericQuery)
}

func BenchmarkDirectMutation_BindingNumeric_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationBindingNumericQuery)
}

func BenchmarkDirectMutation_ObjectProperty_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationObjectPropertyQuery)
}

func BenchmarkDirectMutation_ObjectProperty_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationObjectPropertyQuery)
}

func BenchmarkDirectMutation_ObjectProperty_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationObjectPropertyQuery)
}

func BenchmarkDirectMutation_DynamicKey_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationDynamicKeyQuery)
}

func BenchmarkDirectMutation_DynamicKey_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationDynamicKeyQuery)
}

func BenchmarkDirectMutation_DynamicKey_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationDynamicKeyQuery)
}

func BenchmarkDirectMutation_ArrayIndex_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationArrayIndexQuery)
}

func BenchmarkDirectMutation_ArrayIndex_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationArrayIndexQuery)
}

func BenchmarkDirectMutation_ArrayIndex_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationArrayIndexQuery)
}

func BenchmarkDirectMutation_NestedAugmented_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationNestedAugmentedQuery)
}

func BenchmarkDirectMutation_NestedAugmented_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationNestedAugmentedQuery)
}

func BenchmarkDirectMutation_NestedAugmented_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationNestedAugmentedQuery)
}

func BenchmarkDirectMutation_SafeNoop_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationSafeNoopQuery)
}

func BenchmarkDirectMutation_SafeNoop_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationSafeNoopQuery)
}

func BenchmarkDirectMutation_SafeNoop_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationSafeNoopQuery)
}

func BenchmarkDirectMutation_SafePresent_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationSafePresentQuery)
}

func BenchmarkDirectMutation_SafePresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationSafePresentQuery)
}

func BenchmarkDirectMutation_SafePresent_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationSafePresentQuery)
}

func BenchmarkDirectMutation_UDFCapturedRoot_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationUDFCapturedRootQuery)
}

func BenchmarkDirectMutation_UDFCapturedRoot_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationUDFCapturedRootQuery)
}

func BenchmarkDirectMutation_UDFCapturedRoot_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationUDFCapturedRootQuery)
}

func BenchmarkDirectMutation_DeleteProperty_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationDeletePropertyQuery)
}

func BenchmarkDirectMutation_DeleteProperty_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationDeletePropertyQuery)
}

func BenchmarkDirectMutation_DeleteProperty_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationDeletePropertyQuery)
}

func BenchmarkDirectMutation_DeleteDynamicKey_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationDeleteDynamicKeyQuery)
}

func BenchmarkDirectMutation_DeleteDynamicKey_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationDeleteDynamicKeyQuery)
}

func BenchmarkDirectMutation_DeleteDynamicKey_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationDeleteDynamicKeyQuery)
}

func BenchmarkDirectMutation_DeleteArrayIndex_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationDeleteArrayIndexQuery)
}

func BenchmarkDirectMutation_DeleteArrayIndex_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationDeleteArrayIndexQuery)
}

func BenchmarkDirectMutation_DeleteArrayIndex_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationDeleteArrayIndexQuery)
}

func BenchmarkDirectMutation_DeleteSafeNoop_None(b *testing.B) {
	RunBenchmarkNone(b, directMutationDeleteSafeNoopQuery)
}

func BenchmarkDirectMutation_DeleteSafeNoop_Basic(b *testing.B) {
	RunBenchmarkBasic(b, directMutationDeleteSafeNoopQuery)
}

func BenchmarkDirectMutation_DeleteSafeNoop_Full(b *testing.B) {
	RunBenchmarkFull(b, directMutationDeleteSafeNoopQuery)
}
