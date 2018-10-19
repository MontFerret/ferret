package literals

import (
	"context"
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

func NewObjectLiteralWith(props ...*ObjectPropertyAssignment) *ObjectLiteral {
	return &ObjectLiteral{props}
}

func (l *ObjectLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	obj := values.NewObject()

	for _, el := range l.properties {
		name, err := el.name.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		val, err := el.value.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		if name.Type() != core.StringType {
			return values.None, core.TypeError(name.Type(), core.StringType)
		}

		obj.Set(name.(values.String), val)
	}

	return obj, nil
}
