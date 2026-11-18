package vm

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type debugFields interface {
	runtime.Map
	DebugFieldsPageable() bool
	ReadDebugFields(int, int, func(string, runtime.Value) error) error
}

// Only these concrete objects supply safe native member access.
func debugObject(value runtime.Value) debugFields {
	switch v := value.(type) {
	case *runtime.Object:
		if v != nil {
			return v
		}
	case *data.FastObject:
		if v != nil {
			return v
		}
	}
	return nil
}

func readObjectPage(ctx context.Context, object debugFields, page debugger.VariablePage) ([]DebugChild, error) {
	if !object.DebugFieldsPageable() {
		return nil, errors.ErrUnsupported
	}

	length, _ := object.Length(ctx)
	if page.Filter == debugger.VariableFilterIndexed || page.Start >= int64(length) {
		return nil, nil
	}

	end := min(page.Start+page.Count, int64(length))
	children := make([]DebugChild, 0, end-page.Start)
	err := object.ReadDebugFields(int(page.Start), int(end), func(name string, value runtime.Value) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if !utf8.ValidString(name) {
			return errors.New("invalid UTF-8 debugger property name")
		}

		children = append(children, DebugChild{Kind: debugger.VariableFilterNamed, Name: name, Value: value, Position: page.Start + int64(len(children))})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return children, nil
}

func readBufferedPage(ctx context.Context, page debugger.VariablePage, length int64, at func(context.Context, runtime.Int) (runtime.Value, bool, error)) ([]DebugChild, error) {
	if page.Filter == debugger.VariableFilterNamed || page.Start >= length {
		return nil, nil
	}

	end := min(page.Start+page.Count, length)
	children := make([]DebugChild, 0, end-page.Start)
	// Keep the page's decimal names in one immutable backing allocation.
	var names strings.Builder
	var digits [20]byte
	names.Grow(int(end-page.Start) * len(strconv.AppendInt(digits[:0], end-1, 10)))

	for i := page.Start; i < end; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		value, found, err := at(ctx, runtime.Int(i))
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, errors.New("buffered debugger child disappeared")
		}

		offset := names.Len()
		names.Write(strconv.AppendInt(digits[:0], i, 10))
		children = append(children, DebugChild{Kind: debugger.VariableFilterIndexed, Name: names.String()[offset:], Position: i, Value: value})
	}

	return children, nil
}
