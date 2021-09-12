package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/rs/zerolog"
)

// PAGINATION creates an iterator that goes through pages using CSS selector.
// The iterator starts from the current page i.e. it does not change the page on 1st iteration.
// That allows you to keep scraping logic inside FOR loop.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector for a pagination on the page.
func Pagination(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	logger := logging.
		WithName(logging.FromContext(ctx).With(), "stdlib_html_pagination").
		Str("selector", selector.String()).
		Logger()

	return &Paging{logger, page, selector}, nil
}

var PagingType = core.NewType("paging")

type (
	Paging struct {
		logger   zerolog.Logger
		page     drivers.HTMLPage
		selector drivers.QuerySelector
	}

	PagingIterator struct {
		logger   zerolog.Logger
		page     drivers.HTMLPage
		selector drivers.QuerySelector
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
	return &PagingIterator{p.logger, p.page, p.selector, -1}, nil
}

func (i *PagingIterator) Next(ctx context.Context) (core.Value, core.Value, error) {
	i.pos++

	i.logger.Trace().Int("position", int(i.pos)).Msg("starting to advance iteration")

	if i.pos == 0 {
		i.logger.Trace().Msg("starting point of pagination. nothing to do. exit")
		return values.ZeroInt, values.ZeroInt, nil
	}

	i.logger.Trace().Msg("checking if an element exists...")
	exists, err := i.page.GetMainFrame().ExistsBySelector(ctx, i.selector)

	if err != nil {
		i.logger.Trace().Err(err).Msg("failed to check")

		return values.None, values.None, err
	}

	if !exists {
		i.logger.Trace().Bool("exists", bool(exists)).Msg("element does not exist. exit")

		return values.None, values.None, core.ErrNoMoreData
	}

	i.logger.Trace().Bool("exists", bool(exists)).Msg("element exists. clicking...")

	err = i.page.GetMainFrame().GetElement().ClickBySelector(ctx, i.selector, 1)

	if err != nil {
		i.logger.Trace().Err(err).Msg("failed to click. exit")

		return values.None, values.None, err
	}

	i.logger.Trace().Msg("successfully clicked on element. iteration has succeeded")

	// terminate
	return i.pos, i.pos, nil
}
