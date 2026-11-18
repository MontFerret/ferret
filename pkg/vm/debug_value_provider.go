package vm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Catch panics only at the embedder boundary. Never format a host panic value.
func readHostChild(ctx context.Context, provider runtime.DebugChildrenProvider, kind debugger.VariableFilter, pos int64) (name string, value runtime.Value, found bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("debug child provider panicked")
		}
	}()

	return provider.DebugChildAt(ctx, kind, pos)
}

func readProviderPage(ctx context.Context, provider runtime.DebugChildrenProvider, page debugger.VariablePage) ([]DebugChild, error) {
	info := safeChildrenInfo(provider)
	if page.Filter == debugger.VariableFilterIndexed {
		return appendProviderRange(ctx, provider, info.Indexed, page, nil)
	}

	if page.Filter == debugger.VariableFilterNamed || info.Indexed.Known && info.Indexed.Count == 0 {
		page.Filter = debugger.VariableFilterNamed

		return appendProviderRange(ctx, provider, info.Named, page, nil)
	}

	if !info.Named.Known {
		return nil, errors.ErrUnsupported
	}

	// Subtract the named prefix. Adding both counts can overflow a valid range.
	named := min(page.Count, max(0, info.Named.Count-page.Start))
	var children []DebugChild
	if named > 0 {
		var err error
		children, err = appendProviderRange(ctx, provider, info.Named, debugger.VariablePage{Filter: debugger.VariableFilterNamed, Start: page.Start, Count: named}, nil)
		if err != nil {
			return nil, err
		}
	}

	if named == page.Count {
		return children, nil
	}

	indexed := debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: max(0, page.Start-info.Named.Count), Count: page.Count - named}

	return appendProviderRange(ctx, provider, info.Indexed, indexed, children)
}

func appendProviderRange(ctx context.Context, provider runtime.DebugChildrenProvider, view runtime.DebugChildView, page debugger.VariablePage, children []DebugChild) ([]DebugChild, error) {
	kind, start, count := page.Filter, page.Start, page.Count
	if view.Known {
		if start >= view.Count {
			return children, nil
		}

		count = min(count, view.Count-start)
	}

	if !view.Readable {
		return nil, errors.ErrUnsupported
	}

	for i := start; i < start+count; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		name, value, found, err := readHostChild(ctx, provider, kind, i)
		if err != nil {
			return nil, err
		}

		if !found {
			if view.Known {
				return nil, fmt.Errorf("debug child missing at %s position %d", kind, i)
			}

			break
		}

		if kind == debugger.VariableFilterIndexed {
			name = strconv.FormatInt(i, 10)
		} else if !utf8.ValidString(name) {
			return nil, errors.New("invalid UTF-8 debugger property name")
		}

		children = append(children, DebugChild{Kind: kind, Position: i, Name: name, Value: value})
	}

	return children, nil
}
