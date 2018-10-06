package collections

import (
	"context"
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

	IterableExpression interface {
		core.Expression
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}

	SliceIterator struct {
		values []core.Value
		pos    int
	}

	MapIterator struct {
		values map[string]core.Value
		keys   []string
		pos    int
	}

	ArrayIterator struct {
		values *values.Array
		pos    int
	}

	ObjectIterator struct {
		values *values.Object
		keys   []string
		pos    int
	}

	HTMLNodeIterator struct {
		values values.HTMLNode
		pos    int
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

func NewSliceIterator(input []core.Value) *SliceIterator {
	return &SliceIterator{input, 0}
}

func (iterator *SliceIterator) HasNext() bool {
	return len(iterator.values) > iterator.pos
}

func (iterator *SliceIterator) Next() (core.Value, core.Value, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return val, values.NewInt(idx), nil
	}

	return values.None, values.None, ErrExhausted
}

func NewMapIterator(input map[string]core.Value) *MapIterator {
	return &MapIterator{input, nil, 0}
}

func (iterator *MapIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		keys := make([]string, len(iterator.values))

		i := 0
		for k := range iterator.values {
			keys[i] = k
			i++
		}

		iterator.keys = keys
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *MapIterator) Next() (core.Value, core.Value, error) {
	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val := iterator.values[key]
		iterator.pos++

		return val, values.NewString(key), nil
	}

	return values.None, values.None, ErrExhausted
}

func NewArrayIterator(input *values.Array) *ArrayIterator {
	return &ArrayIterator{input, 0}
}

func (iterator *ArrayIterator) HasNext() bool {
	return int(iterator.values.Length()) > iterator.pos
}

func (iterator *ArrayIterator) Next() (core.Value, core.Value, error) {
	if int(iterator.values.Length()) > iterator.pos {
		idx := iterator.pos
		val := iterator.values.Get(values.NewInt(idx))
		iterator.pos++

		return val, values.NewInt(idx), nil
	}

	return values.None, values.None, ErrExhausted
}

func NewObjectIterator(input *values.Object) *ObjectIterator {
	return &ObjectIterator{input, nil, 0}
}

func (iterator *ObjectIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		iterator.keys = iterator.values.Keys()
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *ObjectIterator) Next() (core.Value, core.Value, error) {
	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val, _ := iterator.values.Get(values.NewString(key))
		iterator.pos++

		return val, values.NewString(key), nil
	}

	return values.None, values.None, ErrExhausted
}

func NewHTMLNodeIterator(input values.HTMLNode) *HTMLNodeIterator {
	return &HTMLNodeIterator{input, 0}
}

func (iterator *HTMLNodeIterator) HasNext() bool {
	return iterator.values.Length() > values.NewInt(iterator.pos)
}

func (iterator *HTMLNodeIterator) Next() (core.Value, core.Value, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := iterator.pos
		val := iterator.values.GetChildNode(values.NewInt(idx))

		iterator.pos++

		return val, values.NewInt(idx), nil
	}

	return values.None, values.None, ErrExhausted
}
