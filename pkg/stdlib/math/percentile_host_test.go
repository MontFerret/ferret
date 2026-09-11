package math_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	percentileCopyList struct {
		runtime.List
		copied runtime.Value
		copies int
	}

	percentileCopySortable struct {
		runtime.Value
		sorts int
	}
)

func (list *percentileCopyList) Copy() runtime.Value {
	list.copies++

	return list.copied
}

func (value *percentileCopySortable) SortAsc(context.Context) error {
	value.sorts++

	return nil
}

func (value *percentileCopySortable) SortDesc(context.Context) error {
	value.sorts++

	return nil
}
