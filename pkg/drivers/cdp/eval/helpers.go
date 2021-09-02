package eval

import (
	"strconv"
	"strings"

	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func CastToValue(input interface{}) (core.Value, error) {
	value, ok := input.(core.Value)

	if !ok {
		return values.None, core.Error(core.ErrInvalidType, "eval return type")
	}

	return value, nil
}

func CastToReference(input interface{}) (runtime.RemoteObject, error) {
	value, ok := input.(runtime.RemoteObject)

	if !ok {
		return runtime.RemoteObject{}, core.Error(core.ErrInvalidType, "eval return type")
	}

	return value, nil
}

func wrapExp(exp string) string {
	return "function () {" + exp + "}"
}

func Unmarshal(obj *runtime.RemoteObject) (core.Value, error) {
	if obj == nil {
		return values.None, nil
	}

	switch obj.Type {
	case "string":
		str, err := strconv.Unquote(string(obj.Value))

		if err != nil {
			return values.None, err
		}

		return values.NewString(str), nil
	case "undefined", "null":
		return values.None, nil
	default:
		return values.Unmarshal(obj.Value)
	}
}

func parseRuntimeException(details *runtime.ExceptionDetails) error {
	if details == nil || details.Exception == nil {
		return nil
	}

	desc := *details.Exception.Description

	if strings.Contains(desc, drivers.ErrNotFound.Error()) {
		return drivers.ErrNotFound
	}

	return core.Error(
		core.ErrUnexpected,
		desc,
	)
}
