package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"strconv"
)

type (
	Iterator interface {
		HasNext() bool
		Next() (ResultSet, error)
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
		set, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		res = append(res, set[0])
	}

	return res, nil
}

func ToMap(iterator Iterator) (map[string]core.Value, error) {
	res := make(map[string]core.Value)

	counter := 0

	for iterator.HasNext() {
		set, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		if len(set) == 1 {
			res[strconv.Itoa(counter)] = set[0]
		} else {
			res[set[1].String()] = set[0]
		}

		counter++
	}

	return res, nil
}

func ToArray(iterator Iterator) (*values.Array, error) {
	res := values.NewArray(10)

	for iterator.HasNext() {
		set, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		res.Push(set[0])
	}

	return res, nil
}

func ToSliceResultSet(iterator Iterator) ([]ResultSet, error) {
	res := make([]ResultSet, 0, 10)

	for iterator.HasNext() {
		set, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		res = append(res, set)
	}

	return res, nil
}
