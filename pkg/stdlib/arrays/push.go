package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// push create a new array with appended value.
// @param array {Any[]} Source array.
// @param value {Any} Target value.
// @param unique {Boolean} Read indicating whether to do uniqueness check.
// @return {Any[]} A new array with appended value.
func Push(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return Append(ctx, args...)
}

// push create a new array with appended value.
// @param array {Any[]} Source array.
// @param value {Any} Target value.
// @return {Any[]} A new array with appended value.
func push2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return append2(ctx, arg1, arg2)
}

// push create a new array with appended value.
// @param array {Any[]} Source array.
// @param value {Any} Target value.
// @param unique {Boolean} Read indicating whether to do uniqueness check.
// @return {Any[]} A new array with appended value.
func push3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return append3(ctx, arg1, arg2, arg3)
}
