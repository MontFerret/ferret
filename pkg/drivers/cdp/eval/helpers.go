package eval

import (
	"bytes"
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

func wrapExp(exp string, args int) string {
	if args == 0 {
		return "function () {\n" + exp + "\n}"
	}

	var buf bytes.Buffer
	lastIndex := args - 1

	for i := 0; i <= args; i++ {
		buf.WriteString("arg")
		buf.WriteString(strconv.Itoa(i))

		if i != lastIndex {
			buf.WriteString(",")
		}
	}

	return "function (" + buf.String() + ") {\n" + exp + "\n}"
}

func Unmarshal(obj runtime.RemoteObject) (core.Value, error) {
	switch obj.Type {
	case "string":
		str, err := strconv.Unquote(string(obj.Value))

		if err != nil {
			return values.None, err
		}

		return values.NewString(str), nil
	case "undefined", "null":
		return values.None, nil
	case "object":
		if obj.Subtype != nil {
			subtype := *obj.Subtype

			if subtype == "null" || subtype == "undefined" {
				return values.None, nil
			}
		}

		return values.Unmarshal(obj.Value)
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
