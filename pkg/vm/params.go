package vm

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var emptyParams = make([]runtime.Value, 0)

func bindParams(required []string, env *Environment) ([]runtime.Value, error) {
	if len(required) == 0 {
		return emptyParams, nil
	}

	paramsSlots := make([]runtime.Value, len(required))
	var missedParams []string

	for idx, name := range required {
		val, exists := env.Params[name]
		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(required))
			}

			missedParams = append(missedParams, "@"+name)
			continue
		}

		paramsSlots[idx] = val
	}

	if len(missedParams) > 0 {
		return nil, runtime.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return paramsSlots, nil
}
