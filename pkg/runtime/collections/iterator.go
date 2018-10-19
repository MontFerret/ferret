package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Iterator interface {
		HasNext() bool
		Next() (value core.Value, key core.Value, err error)
	}

	Iterable interface {
		Iterate() Iterator
	}
)

func ToIterator(value core.Value) (Iterator, error) {
	switch value.Type() {
	case core.ArrayType:
		return NewArrayIterator(value.(*values.Array)), nil
	case core.ObjectType:
		return NewObjectIterator(value.(*values.Object)), nil
	case core.HTMLElementType, core.HTMLDocumentType:
		return NewHTMLNodeIterator(value.(values.HTMLNode)), nil
	default:
		return nil, core.TypeError(
			value.Type(),
			core.ArrayType,
			core.ObjectType,
			core.HTMLDocumentType,
			core.HTMLElementType,
		)
	}
}

func ToSlice(iterator Iterator) ([]core.Value, error) {
	res := make([]core.Value, 0, 10)

	for iterator.HasNext() {
		item, _, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func ToMap(iterator Iterator) (map[string]core.Value, error) {
	res := make(map[string]core.Value)

	for iterator.HasNext() {
		item, key, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		res[key.String()] = item
	}

	return res, nil
}

func ToArray(iterator Iterator) (*values.Array, error) {
	res := values.NewArray(10)

	for iterator.HasNext() {
		item, _, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		res.Push(item)
	}

	return res, nil
}
