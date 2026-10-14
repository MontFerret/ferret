package msgpack

import (
	"context"
	"errors"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type encoder struct {
	pre  []encoding.PreEncoderHook
	post []encoding.PostEncoderHook
}

func (enc encoder) encodeValue(ctx context.Context, op *codecutil.Operation, menc *vmmsgpack.Encoder, value runtime.Value) error {
	op.Adopt(value, len(enc.pre) > 0 || len(enc.post) > 0)
	if err := op.Step(); err != nil {
		op.Adopt(value, true)
		return err
	}

	if err := enc.runPreHooks(ctx, value); err != nil {
		op.Adopt(value, true)
		return err
	}

	if value == nil || value == runtime.None {
		err := menc.EncodeNil()

		if hookErr := enc.runPostHooks(ctx, runtime.None, err); hookErr != nil {
			return errors.Join(err, hookErr)
		}

		return err
	}

	var err error

	if supportsCustomEncoding(value) {
		err = enc.encodeAny(op, menc, value)
	} else {
		switch v := value.(type) {
		case runtime.Boolean:
			err = menc.EncodeBool(bool(v))
		case runtime.Int:
			err = menc.EncodeInt(int64(v))
		case runtime.Float:
			err = menc.EncodeFloat64(float64(v))
		case runtime.Duration:
			err = menc.EncodeString(v.String())
		case runtime.String:
			err = menc.EncodeString(string(v))
		case *runtime.Regexp:
			err = menc.EncodeString(v.String())
		case runtime.Binary:
			err = menc.EncodeBytes(v)
		case runtime.DateTime:
			err = menc.EncodeTime(v.Time)
		case *runtime.Range:
			err = enc.encodeRange(ctx, op, menc, v)
		case runtime.Map:
			err = enc.encodeMap(ctx, op, menc, v)
		case runtime.List:
			err = enc.encodeList(ctx, op, menc, v)
		case runtime.Iterable:
			err = enc.encodeIterable(ctx, op, menc, v)
		case runtime.Unwrappable:
			err = enc.encodeAny(op, menc, v.Unwrap())
		default:
			err = enc.encodeAny(op, menc, v)
		}
	}

	if len(enc.post) > 0 {
		op.Adopt(value, true)
	}

	if hookErr := enc.runPostHooks(ctx, value, err); hookErr != nil {
		err = errors.Join(err, hookErr)
	}

	if err != nil || len(enc.pre) > 0 || len(enc.post) > 0 {
		op.Adopt(value, true)
	}

	return err
}

func (enc encoder) encodeAny(op *codecutil.Operation, menc *vmmsgpack.Encoder, value any) error {
	if err := op.Check(); err != nil {
		return err
	}

	err := menc.Encode(value)
	// A custom encoder may ignore its writer's error or return another cause.
	// Preserve the codec-owned writer failure in either case.
	if writer, ok := menc.Writer().(*codecutil.Writer); ok && writer.Error() != nil && !errors.Is(err, writer.Error()) {
		err = errors.Join(err, writer.Error())
	}

	if err != nil {
		return err
	}

	return op.Check()
}

func (enc encoder) encodeMap(ctx context.Context, op *codecutil.Operation, menc *vmmsgpack.Encoder, value runtime.Map) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	nativeLength, err := collectionLength(length)
	if err != nil {
		return err
	}

	if err := menc.EncodeMapLen(nativeLength); err != nil {
		return err
	}

	_, stored := value.(*runtime.Object)
	return value.ForEach(ctx, func(ctx context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		op.Adopt(val, !stored)
		op.Adopt(key, !stored)

		if err := menc.EncodeString(key.String()); err != nil {
			op.Adopt(val, true)
			op.Adopt(key, true)

			return false, err
		}

		if err := enc.encodeValue(ctx, op, menc, val); err != nil {
			return false, err
		}

		return true, nil
	})
}

func (enc encoder) encodeList(ctx context.Context, op *codecutil.Operation, menc *vmmsgpack.Encoder, value runtime.List) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	nativeLength, err := collectionLength(length)
	if err != nil {
		return err
	}

	if err := menc.EncodeArrayLen(nativeLength); err != nil {
		return err
	}

	_, stored := value.(*runtime.Array)
	for i := runtime.Int(0); i < length; i++ {
		item, err := value.At(ctx, i)
		if err != nil {
			return err
		}

		// Native storage immediately enters encodeValue, which registers before
		// its first checkpoint. Host acquisitions also discover stored children.
		if !stored {
			op.Adopt(item, true)
		}

		if err := enc.encodeValue(ctx, op, menc, item); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) encodeIterable(ctx context.Context, op *codecutil.Operation, menc *vmmsgpack.Encoder, value runtime.Iterable) (err error) {
	iter, err := codecutil.NewIterator(ctx, value)
	if err != nil {
		return err
	}

	list := runtime.NewArray(0)
	defer func() {
		if closeErr := iter.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		if err != nil || len(enc.pre) > 0 || len(enc.post) > 0 {
			// Hooks or host traversal can mutate another collected container.
			// Discover its stored children before leaving the internal buffer.
			op.Adopt(list, true)
		}
	}()
	err = runtime.ForEachIter(ctx, iter.Iterator, func(ctx context.Context, item, key runtime.Value) (runtime.Boolean, error) {
		// Yielding transfers discovery to the result before a checkpoint can stop
		// collection. This callback is separate from user encoding hooks.
		op.Adopt(item, true)
		op.Adopt(key, true)

		if err := op.Step(); err != nil {
			return false, err
		}

		if err := list.Append(ctx, item); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	if len(enc.pre) > 0 || len(enc.post) > 0 {
		op.Adopt(list, true)
	}

	return enc.encodeList(ctx, op, menc, list)
}

func (enc encoder) encodeRange(ctx context.Context, op *codecutil.Operation, menc *vmmsgpack.Encoder, value *runtime.Range) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	nativeLength, err := collectionLength(length)
	if err != nil {
		return err
	}

	if err := menc.EncodeArrayLen(nativeLength); err != nil {
		return err
	}

	start := value.Start()
	ascending := start <= value.End()

	for offset := runtime.Int(0); offset < length; offset++ {
		item := start + offset
		if !ascending {
			item = start - offset
		}

		if err := enc.encodeValue(ctx, op, menc, item); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPreHooks(ctx context.Context, value runtime.Value) error {
	if len(enc.pre) == 0 {
		return nil
	}

	for _, hook := range enc.pre {
		if err := hook(ctx, value); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPostHooks(ctx context.Context, value runtime.Value, err error) error {
	if len(enc.post) == 0 {
		return nil
	}

	for _, hook := range enc.post {
		if hookErr := hook(ctx, value, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
