package eval

import (
	"fmt"

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

	if obj.Type != "undefined" {
		return values.Unmarshal(obj.Value)
	}

	return values.None, nil
}
