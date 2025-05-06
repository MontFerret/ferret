package eval

import (
	"context"
	"strconv"

	runtime2 "github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	ValueLoader interface {
		Load(
			ctx context.Context,
			frameID page.FrameID,
			remoteType RemoteObjectType,
			remoteClass RemoteClassName,
			id runtime.RemoteObjectID,
		) (core.Value, error)
	}

	ValueLoaderFn func(
		ctx context.Context,
		frameID page.FrameID,
		remoteType RemoteObjectType,
		remoteClass RemoteClassName,
		id runtime.RemoteObjectID,
	) (core.Value, error)

	Resolver struct {
		runtime cdp.Runtime
		frameID page.FrameID
		loader  ValueLoader
	}
)

func (f ValueLoaderFn) Load(
	ctx context.Context,
	frameID page.FrameID,
	remoteType RemoteObjectType,
	remoteClass RemoteClassName,
	id runtime.RemoteObjectID,
) (core.Value, error) {
	return f(ctx, frameID, remoteType, remoteClass, id)
}

func NewResolver(runtime cdp.Runtime, frameID page.FrameID) *Resolver {
	return &Resolver{runtime, frameID, nil}
}

func (r *Resolver) SetLoader(loader ValueLoader) *Resolver {
	r.loader = loader

	return r
}

func (r *Resolver) ToValue(ctx context.Context, ref runtime.RemoteObject) (core.Value, error) {
	// It's not an actual ref but rather a plain value
	if ref.ObjectID == nil {
		if ref.Value != nil {
			return internal.Unmarshal(ref.Value)
		}

		return core.None, nil
	}

	subtype := ToRemoteObjectType(ref)

	switch subtype {
	case NullObjectType, UndefinedObjectType:
		return core.None, nil
	case ArrayObjectType:
		props, err := r.runtime.GetProperties(ctx, runtime.NewGetPropertiesArgs(*ref.ObjectID).SetOwnProperties(true))

		if err != nil {
			return core.None, err
		}

		if props.ExceptionDetails != nil {
			exception := *props.ExceptionDetails

			return core.None, errors.New(exception.Text)
		}

		result := runtime2.NewArray(len(props.Result))

		for _, descr := range props.Result {
			if !descr.Enumerable {
				continue
			}

			if descr.Value == nil {
				continue
			}

			el, err := r.ToValue(ctx, *descr.Value)

			if err != nil {
				return core.None, err
			}

			result.Push(el)
		}

		return result, nil
	case NodeObjectType:
		// is it even possible?
		if ref.ObjectID == nil {
			return internal.Unmarshal(ref.Value)
		}

		return r.loadValue(ctx, NodeObjectType, ToRemoteClassName(ref), *ref.ObjectID)
	default:
		switch ToRemoteType(ref) {
		case StringType:
			str, err := strconv.Unquote(string(ref.Value))

			if err != nil {
				return core.None, err
			}

			return core.NewString(str), nil
		case ObjectType:
			if subtype == NullObjectType || subtype == UnknownObjectType {
				return core.None, nil
			}

			return internal.Unmarshal(ref.Value)
		default:
			return internal.Unmarshal(ref.Value)
		}
	}
}

func (r *Resolver) ToElement(ctx context.Context, ref runtime.RemoteObject) (drivers.HTMLElement, error) {
	if ref.ObjectID == nil {
		return nil, core.Error(core.ErrInvalidArgument, "ref id")
	}

	val, err := r.loadValue(ctx, ToRemoteObjectType(ref), ToRemoteClassName(ref), *ref.ObjectID)

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
		return core.None, err
	}

	if err := parseRuntimeException(res.ExceptionDetails); err != nil {
		return core.None, err
	}

	for _, prop := range res.Result {
		if prop.Name == propName {
			if prop.Value != nil {
				return r.ToValue(ctx, *prop.Value)
			}

			return core.None, nil
		}
	}

	return core.None, nil
}

func (r *Resolver) ToProperties(
	ctx context.Context,
	id runtime.RemoteObjectID,
) (*runtime2.Array, error) {
	res, err := r.runtime.GetProperties(
		ctx,
		runtime.NewGetPropertiesArgs(id),
	)

	if err != nil {
		return runtime2.EmptyArray(), err
	}

	if err := parseRuntimeException(res.ExceptionDetails); err != nil {
		return runtime2.EmptyArray(), err
	}

	arr := runtime2.NewArray(len(res.Result))

	for _, prop := range res.Result {
		if prop.Value != nil {
			val, err := r.ToValue(ctx, *prop.Value)

			if err != nil {
				return runtime2.EmptyArray(), err
			}

			arr.Push(val)
		}
	}

	return arr, nil
}

func (r *Resolver) loadValue(ctx context.Context, remoteType RemoteObjectType, remoteClass RemoteClassName, id runtime.RemoteObjectID) (core.Value, error) {
	if r.loader == nil {
		return core.None, core.Error(core.ErrNotImplemented, "ValueLoader")
	}

	return r.loader.Load(ctx, r.frameID, remoteType, remoteClass, id)
}
