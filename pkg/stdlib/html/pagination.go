package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Pagination creates an iterator that goes through pages using CSS selector.
// The iterator starts from the current page i.e. it does not change the page on 1st iteration.
// That allows you to keep scraping logic inside FOR loop.
// @param doc (Document) - Target document.
// @param selector (String) - CSS selector for a pagination on the page.
func Pagination(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return &Paging{doc, selector}, nil
}

type (
	Paging struct {
		document *dynamic.HTMLDocument
		selector values.String
	}

	PagingIterator struct {
		document *dynamic.HTMLDocument
		selector values.String
		pos      values.Int
	}
)

func (p *Paging) MarshalJSON() ([]byte, error) {
	return nil, core.ErrInvalidOperation
}

func (p *Paging) Type() core.Type {
	return core.CustomType
}

func (p *Paging) String() string {
	return core.CustomType.String()
}

func (p *Paging) Compare(_ core.Value) int {
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

func (p *Paging) Iterate(_ context.Context) (collections.CollectionIterator, error) {
	return &PagingIterator{p.document, p.selector, -1}, nil
}

func (i *PagingIterator) Next(_ context.Context) (core.Value, core.Value, error) {
	i.pos++

	if i.pos == 0 {
		return values.ZeroInt, values.ZeroInt, nil
	}

	clicked, err := i.document.ClickBySelector(i.selector)

	if err != nil {
		return values.None, values.None, err
	}

	if clicked {
		return i.pos, i.pos, nil
	}

	// terminate
	return values.None, values.None, nil
}
