package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("first", First).
		Add("flatten", flatten1).
		Add("last", Last).
		Add("pop", Pop).
		Add("shift", Shift).
		Add("sorted", Sorted).
		Add("sorted_unique", SortedUnique).
		Add("unique", Unique)

	ns.Function().A2().
		Add("append", append2).
		Add("flatten", flatten2).
		Add("nth", Nth).
		Add("position", position2).
		Add("push", push2).
		Add("remove_value", removeValue2).
		Add("remove_nth", RemoveNth).
		Add("remove_values", RemoveValues).
		Add("slice", slice2).
		Add("unshift", unshift2)

	ns.Function().A3().
		Add("append", append3).
		Add("position", position3).
		Add("push", push3).
		Add("remove_value", removeValue3).
		Add("slice", slice3).
		Add("unshift", unshift3)

	ns.Function().Var().
		Add("intersection", Intersection).
		Add("minus", Minus).
		Add("outersection", Outersection).
		Add("union", Union).
		Add("union_distinct", UnionDistinct)
}

func ToUniqueList(ctx context.Context, list runtime.List) (runtime.List, error) {
	seen := valueset.New(0)

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		added, err := seen.Add(ctx, value)
		return runtime.Boolean(added), err
	})
}
