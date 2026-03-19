package msgpack

import (
	"context"
	"errors"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type encoder struct {
	pre  []encoding.PreEncoderHook
	post []encoding.PostEncoderHook
}

func (enc encoder) encodeValue(ctx context.Context, menc *vmmsgpack.Encoder, value runtime.Value) error {
	if err := enc.runPreHooks(value); err != nil {
		return err
	}

	if value == nil || value == runtime.None {
		err := menc.EncodeNil()

		if hookErr := enc.runPostHooks(runtime.None, err); hookErr != nil {
			return errors.Join(err, hookErr)
		}

		return err
	}

	var err error

	if supportsCustomEncoding(value) {
		err = enc.encodeAny(menc, value)
	} else {
		switch v := value.(type) {
		case runtime.Boolean:
			err = menc.EncodeBool(bool(v))
		case runtime.Int:
			err = menc.EncodeInt(int64(v))
		case runtime.Float:
			err = menc.EncodeFloat64(float64(v))
		case runtime.String:
			err = menc.EncodeString(string(v))
		case runtime.Binary:
			err = menc.EncodeBytes(v)
		case runtime.DateTime:
			err = menc.EncodeTime(v.Time)
		case runtime.Map:
			err = enc.encodeMap(ctx, menc, v)
		case runtime.List:
			err = enc.encodeList(ctx, menc, v)
		case runtime.Iterable:
			list, e := runtime.ToList(ctx, value)
			if e != nil {
				err = e
				break
			}

			err = enc.encodeList(ctx, menc, list)
		case runtime.Unwrappable:
			err = enc.encodeAny(menc, v.Unwrap())
		default:
			err = enc.encodeAny(menc, v)
		}
	}

	if hookErr := enc.runPostHooks(value, err); hookErr != nil {
		return errors.Join(err, hookErr)
	}

	return err
}

func supportsCustomEncoding(value runtime.Value) bool {
	switch value.(type) {
	case vmmsgpack.CustomEncoder, vmmsgpack.Marshaler:
		return true
	default:
		return false
	}
}

func (enc encoder) encodeAny(menc *vmmsgpack.Encoder, value any) error {
	return menc.Encode(value)
}

func (enc encoder) encodeMap(ctx context.Context, menc *vmmsgpack.Encoder, value runtime.Map) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := menc.EncodeMapLen(int(length)); err != nil {
		return err
	}

	return value.ForEach(ctx, func(_ context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		if err := menc.EncodeString(key.String()); err != nil {
			return false, err
		}

		if err := enc.encodeValue(ctx, menc, val); err != nil {
			return false, err
		}

		return true, nil
	})
}

func (enc encoder) encodeList(ctx context.Context, menc *vmmsgpack.Encoder, value runtime.List) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := menc.EncodeArrayLen(int(length)); err != nil {
		return err
	}

	for i := runtime.Int(0); i < length; i++ {
		item, err := value.At(ctx, i)
		if err != nil {
			return err
		}

		if err := enc.encodeValue(ctx, menc, item); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPreHooks(value runtime.Value) error {
	if len(enc.pre) == 0 {
		return nil
	}

	for _, hook := range enc.pre {
		if err := hook(value); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPostHooks(value runtime.Value, err error) error {
	if len(enc.post) == 0 {
		return nil
	}

	for _, hook := range enc.post {
		if hookErr := hook(value, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
