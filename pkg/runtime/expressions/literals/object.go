package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	ObjectPropertyAssignment struct {
		name  core.Expression
		value core.Expression
	}

	ObjectLiteral struct {
		properties []*ObjectPropertyAssignment
	}
)

func NewObjectPropertyAssignment(name, value core.Expression) *ObjectPropertyAssignment {
	return &ObjectPropertyAssignment{name, value}
}

func NewObjectLiteral() *ObjectLiteral {
	return &ObjectLiteral{make([]*ObjectPropertyAssignment, 0, 10)}
}

func NewObjectLiteralWith(props ...*ObjectPropertyAssignment) *ObjectLiteral {
	return &ObjectLiteral{props}
}

func (l *ObjectLiteral) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	obj, err := l.doExec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewObjectIterator(obj), nil
}

func (l *ObjectLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	arr, err := l.doExec(ctx, scope)

	if err != nil {
		return values.None, err
	}

	return arr, nil
}

func (l *ObjectLiteral) doExec(ctx context.Context, scope *core.Scope) (*values.Object, error) {
	obj := values.NewObject()

	for _, el := range l.properties {
		name, err := el.name.Exec(ctx, scope)

		if err != nil {
			return nil, err
		}

		val, err := el.value.Exec(ctx, scope)

		if err != nil {
			return nil, err
		}

		if name.Type() != core.StringType {
			return nil, core.TypeError(name.Type(), core.StringType)
		}

		obj.Set(name.(values.String), val)
	}

	return obj, nil
}
