package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("FIRST", First).
		Add("FLATTEN", flatten1).
		Add("LAST", Last).
		Add("POP", Pop).
		Add("SHIFT", Shift).
		Add("SORTED", Sorted).
		Add("SORTED_UNIQUE", SortedUnique).
		Add("UNIQUE", Unique)

	ns.Function().A2().
		Add("APPEND", append2).
		Add("FLATTEN", flatten2).
		Add("NTH", Nth).
		Add("POSITION", position2).
		Add("PUSH", push2).
		Add("REMOVE_VALUE", removeValue2).
		Add("REMOVE_NTH", RemoveNth).
		Add("REMOVE_VALUES", RemoveValues).
		Add("SLICE", slice2).
		Add("UNSHIFT", unshift2)

	ns.Function().A3().
		Add("APPEND", append3).
		Add("POSITION", position3).
		Add("PUSH", push3).
		Add("REMOVE_VALUE", removeValue3).
		Add("SLICE", slice3).
		Add("UNSHIFT", unshift3)

	ns.Function().Var().
		Add("INTERSECTION", Intersection).
		Add("MINUS", Minus).
		Add("OUTERSECTION", Outersection).
		Add("UNION", Union).
		Add("UNION_DISTINCT", UnionDistinct)
}

func ToUniqueList(ctx context.Context, list runtime.List) (runtime.List, error) {
	seen := valueset.New(0)

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		added, err := seen.Add(ctx, value)
		return runtime.Boolean(added), err
	})
}
