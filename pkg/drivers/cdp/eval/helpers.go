package eval

import (
	"fmt"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

func PrepareEval(exp string) string {
	return fmt.Sprintf("((function () {%s})())", exp)
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
