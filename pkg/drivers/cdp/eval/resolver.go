package eval

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"strconv"
)

type (
	ValueLoader func(ctx context.Context, remoteType RemoteObjectType, id runtime.RemoteObjectID) (core.Value, error)

	Resolver struct {
		runtime cdp.Runtime
		loader  ValueLoader
	}
)

func NewResolver(runtime cdp.Runtime) *Resolver {
	return &Resolver{runtime, nil}
}

func (r *Resolver) SetLoader(loader ValueLoader) *Resolver {
	r.loader = loader

	return r
}

func (r *Resolver) ToValue(ctx context.Context, ref runtime.RemoteObject) (core.Value, error) {
	// It's not an actual ref but rather a plain value
	if ref.ObjectID == nil {
		return values.Unmarshal(ref.Value)
	}

	subtype := ToRemoteObjectType(ref)

	switch subtype {
	case NullObjectType, UndefinedObjectType:
		return values.None, nil
	case ArrayObjectType:
		props, err := r.runtime.GetProperties(ctx, runtime.NewGetPropertiesArgs(*ref.ObjectID).SetOwnProperties(true))

		if err != nil {
			return values.None, err
		}

		if props.ExceptionDetails != nil {
			exception := *props.ExceptionDetails

			return values.None, errors.New(exception.Text)
		}

		result := values.NewArray(len(props.Result))

		for _, descr := range props.Result {
			if !descr.Enumerable {
				continue
			}

			if descr.Value == nil {
				continue
			}

			el, err := r.ToValue(ctx, *descr.Value)

			if err != nil {
				return values.None, err
			}

			result.Push(el)
		}

		return result, nil
	case NodeObjectType:
		// could it be possible?
		if ref.ObjectID == nil {
			return values.Unmarshal(ref.Value)
		}

		return r.loadValue(ctx, NodeObjectType, *ref.ObjectID)
	default:
		switch ToRemoteType(ref) {
		case StringType:
			str, err := strconv.Unquote(string(ref.Value))

			if err != nil {
				return values.None, err
			}

			return values.NewString(str), nil
		case ObjectType:
			if subtype == NullObjectType || subtype == UnknownObjectType {
				return values.None, nil
			}

			return values.Unmarshal(ref.Value)
		default:
			return values.Unmarshal(ref.Value)
		}
	}
}

func (r *Resolver) ToElement(ctx context.Context, ref runtime.RemoteObject) (drivers.HTMLElement, error) {
	if ref.ObjectID == nil {
		return nil, core.Error(core.ErrInvalidArgument, "ref id")
	}

	val, err := r.loadValue(ctx, ToRemoteObjectType(ref), *ref.ObjectID)

	if err != nil {
		return nil, err
	}

	return drivers.ToElement(val)
}

func (r *Resolver) ToProperty(
	ctx context.Context,
	id runtime.RemoteObjectID,
	propName string,
) (core.Value, error) {
	res, err := r.runtime.GetProperties(
		ctx,
		runtime.NewGetPropertiesArgs(id),
	)

	if err != nil {
		return values.None, err
	}

	if err := parseRuntimeException(res.ExceptionDetails); err != nil {
		return values.None, err
	}

	for _, prop := range res.Result {
		if prop.Name == propName {
			if prop.Value != nil {
				return r.ToValue(ctx, *prop.Value)
			}

			return values.None, nil
		}
	}

	return values.None, nil
}

func (r *Resolver) ToProperties(
	ctx context.Context,
	id runtime.RemoteObjectID,
) (*values.Array, error) {
	res, err := r.runtime.GetProperties(
		ctx,
		runtime.NewGetPropertiesArgs(id),
	)

	if err != nil {
		return values.EmptyArray(), err
	}

	if err := parseRuntimeException(res.ExceptionDetails); err != nil {
		return values.EmptyArray(), err
	}

	arr := values.NewArray(len(res.Result))

	for _, prop := range res.Result {
		if prop.Value != nil {
			val, err := r.ToValue(ctx, *prop.Value)

			if err != nil {
				return values.EmptyArray(), err
			}

			arr.Push(val)
		}
	}

	return arr, nil
}

func (r *Resolver) loadValue(ctx context.Context, remoteType RemoteObjectType, id runtime.RemoteObjectID) (core.Value, error) {
	if r.loader == nil {
		return values.None, core.Error(core.ErrNotImplemented, "ValueLoader")
	}

	return r.loader(ctx, remoteType, id)
}
