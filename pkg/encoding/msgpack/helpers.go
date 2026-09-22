package msgpack

import (
	"fmt"
	"math"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func supportsCustomEncoding(value runtime.Value) bool {
	switch value.(type) {
	case vmmsgpack.CustomEncoder, vmmsgpack.Marshaler:
		return true
	default:
		return false
	}
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

func unsignedIntValue(value uint64) (runtime.Value, error) {
	if value <= math.MaxInt64 {
		return runtime.Int(value), nil
	}

	return runtime.None, fmt.Errorf("msgpack: integer %d exceeds runtime range", value)
}
