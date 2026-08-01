package internal

import "github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"

type (
	udfCaptureEnv struct {
		scopes []map[string]captureBindingInfo
	}

	udfCaptureState struct {
		captures map[core.BindingID]core.UDFCapture
		order    []core.BindingID
		outer    map[core.BindingID]captureBindingInfo
		owned    map[core.BindingID]captureBindingInfo
		callees  []*core.UDFInfo
	}
)

func (e *udfCaptureEnv) push() {
	e.scopes = append(e.scopes, make(map[string]captureBindingInfo))
}

func (e *udfCaptureEnv) pop() {
	if len(e.scopes) > 0 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *udfCaptureEnv) addBinding(binding captureBindingInfo) {
	if len(e.scopes) == 0 {
		return
	}

	e.scopes[len(e.scopes)-1][binding.Name] = binding
}

func (e *udfCaptureEnv) bindingsByID() map[core.BindingID]captureBindingInfo {
	out := make(map[core.BindingID]captureBindingInfo)

	for _, scope := range e.scopes {
		for _, binding := range scope {
			if binding.ID == core.InvalidBindingID {
				continue
			}

			out[binding.ID] = binding
		}
	}

	return out
}

func (e *udfCaptureEnv) currentHas(name string) bool {
	if len(e.scopes) == 0 {
		return false
	}

	_, ok := e.scopes[len(e.scopes)-1][name]
	return ok
}

func (e *udfCaptureEnv) resolveBinding(name string) (captureBindingInfo, bool) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if binding, ok := e.scopes[i][name]; ok {
			return binding, true
		}
	}

	return captureBindingInfo{}, false
}
