package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// PAGINATION creates an iterator that goes through pages using CSS selector.
// The iterator starts from the current page i.e. it does not change the page on 1st iteration.
// That allows you to keep scraping logic inside FOR loop.
// @param doc (Open) - Target document.
// @param selector (String) - CSS selector for a pagination on the page.
func Pagination(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return &Paging{page, selector}, nil
}

var PagingType = core.NewType("paging")

type (
	Paging struct {
		page     drivers.HTMLPage
		selector values.String
	}

	PagingIterator struct {
		page     drivers.HTMLPage
		selector values.String
		pos      values.Int
	}
)

func (p *Paging) MarshalJSON() ([]byte, error) {
	return nil, core.ErrInvalidOperation
}

func (p *Paging) Type() core.Type {
	return PagingType
}

func (p *Paging) String() string {
	return PagingType.String()
}

func (p *Paging) Compare(_ core.Value) int64 {
	return 1
}

func (p *Paging) Unwrap() interface{} {
	return nil
}

func (p *Paging) Hash() uint64 {
	return 0
}

func (p *Paging) Copy() core.Value {
	return values.None
}

func (p *Paging) Iterate(_ context.Context) (core.Iterator, error) {
	return &PagingIterator{p.page, p.selector, -1}, nil
}

func (i *PagingIterator) Next(ctx context.Context) (core.Value, core.Value, error) {
	i.pos++

	if i.pos == 0 {
		return values.ZeroInt, values.ZeroInt, nil
	}

	exists, err := i.page.GetMainFrame().ExistsBySelector(ctx, i.selector)

	if err != nil {
		return values.None, values.None, err
	}

	if !exists {
		return values.None, values.None, core.ErrNoMoreData
	}

	err = i.page.GetMainFrame().GetElement().ClickBySelector(ctx, i.selector, 1)

	if err != nil {
		return values.None, values.None, err
	}

	// terminate
	return i.pos, i.pos, nil
}
