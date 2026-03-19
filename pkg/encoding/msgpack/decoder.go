package msgpack

import (
	"context"
	"fmt"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const maxExactIntegerFloat64 = uint64(1 << 53)

var (
	maxRuntimeInt = int64(^uint(0) >> 1)
	minRuntimeInt = -maxRuntimeInt - 1
)

type decoder struct {
	pre  []encoding.PreDecoderHook
	post []encoding.PostDecoderHook
}

func (dec decoder) decodeValue(ctx context.Context, d *vmmsgpack.Decoder) (runtime.Value, error) {
	code, err := d.PeekCode()
	if err != nil {
		return runtime.None, err
	}

	switch {
	case code == msgpcode.Nil:
		if err := d.DecodeNil(); err != nil {
			return runtime.None, err
		}

		return runtime.None, nil
	case code == msgpcode.False || code == msgpcode.True:
		value, err := d.DecodeBool()
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewBoolean(value), nil
	case msgpcode.IsString(code):
		value, err := d.DecodeString()
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewString(value), nil
	case msgpcode.IsBin(code):
		value, err := d.DecodeBytes()
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewBinary(value), nil
	case code == msgpcode.Float || code == msgpcode.Double:
		value, err := d.DecodeFloat64()
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewFloat(value), nil
	case isSignedIntCode(code):
		value, err := d.DecodeInt64()
		if err != nil {
			return runtime.None, err
		}

		return signedIntValue(value)
	case isUnsignedIntCode(code):
		value, err := d.DecodeUint64()
		if err != nil {
			return runtime.None, err
		}

		return unsignedIntValue(value)
	case isArrayCode(code):
		return dec.decodeArray(ctx, d)
	case isMapCode(code):
		return dec.decodeMap(ctx, d)
	case msgpcode.IsExt(code):
		value, err := d.DecodeTime()
		if err != nil {
			return runtime.None, err
		}

		return runtime.NewDateTime(value), nil
	default:
		return runtime.None, fmt.Errorf("msgpack: unsupported code %x", code)
	}
}

func (dec decoder) decodeArray(ctx context.Context, d *vmmsgpack.Decoder) (runtime.Value, error) {
	size, err := d.DecodeArrayLen()
	if err != nil {
		return runtime.None, err
	}

	if size < 0 {
		return runtime.None, nil
	}

	arr := runtime.NewArray(size)
	for i := 0; i < size; i++ {
		value, err := dec.decodeValue(ctx, d)
		if err != nil {
			return runtime.None, err
		}

		if err := arr.Append(ctx, value); err != nil {
			return runtime.None, err
		}
	}

	return arr, nil
}

func (dec decoder) decodeMap(ctx context.Context, d *vmmsgpack.Decoder) (runtime.Value, error) {
	size, err := d.DecodeMapLen()
	if err != nil {
		return runtime.None, err
	}

	if size < 0 {
		return runtime.None, nil
	}

	obj := runtime.NewObjectOf(size)
	for i := 0; i < size; i++ {
		code, err := d.PeekCode()
		if err != nil {
			return runtime.None, err
		}

		var key string
		if msgpcode.IsString(code) {
			key, err = d.DecodeString()
			if err != nil {
				return runtime.None, err
			}
		} else {
			keyValue, err := dec.decodeValue(ctx, d)
			if err != nil {
				return runtime.None, err
			}

			key = keyValue.String()
		}

		value, err := dec.decodeValue(ctx, d)
		if err != nil {
			return runtime.None, err
		}

		if err := obj.Set(ctx, runtime.NewString(key), value); err != nil {
			return runtime.None, err
		}
	}

	return obj, nil
}

func isSignedIntCode(code byte) bool {
	return code >= msgpcode.NegFixedNumLow ||
		code == msgpcode.Int8 ||
		code == msgpcode.Int16 ||
		code == msgpcode.Int32 ||
		code == msgpcode.Int64
}

func isUnsignedIntCode(code byte) bool {
	return code <= msgpcode.PosFixedNumHigh ||
		code == msgpcode.Uint8 ||
		code == msgpcode.Uint16 ||
		code == msgpcode.Uint32 ||
		code == msgpcode.Uint64
}

func isArrayCode(code byte) bool {
	return msgpcode.IsFixedArray(code) || code == msgpcode.Array16 || code == msgpcode.Array32
}

func isMapCode(code byte) bool {
	return msgpcode.IsFixedMap(code) || code == msgpcode.Map16 || code == msgpcode.Map32
}

func signedIntValue(value int64) (runtime.Value, error) {
	if value >= minRuntimeInt && value <= maxRuntimeInt {
		return runtime.Int(value), nil
	}

	floatValue := float64(value)
	if int64(floatValue) == value {
		return runtime.NewFloat(floatValue), nil
	}

	return runtime.None, fmt.Errorf("msgpack: integer %d exceeds runtime range", value)
}

func unsignedIntValue(value uint64) (runtime.Value, error) {
	if value <= uint64(maxRuntimeInt) {
		return runtime.Int(value), nil
	}

	if value <= maxExactIntegerFloat64 {
		return runtime.NewFloat(float64(value)), nil
	}

	return runtime.None, fmt.Errorf("msgpack: integer %d exceeds runtime range", value)
}

func (dec decoder) runPreHooks(data []byte) error {
	if len(dec.pre) == 0 {
		return nil
	}

	for _, hook := range dec.pre {
		if err := hook(data); err != nil {
			return err
		}
	}

	return nil
}

func (dec decoder) runPostHooks(data []byte, err error) error {
	if len(dec.post) == 0 {
		return nil
	}

	for _, hook := range dec.post {
		if hookErr := hook(data, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
