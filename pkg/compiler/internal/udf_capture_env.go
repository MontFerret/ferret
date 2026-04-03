package internal

type (
	udfCaptureEnv struct {
		scopes []map[string]captureBindingInfo
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

func (e *udfCaptureEnv) add(name string) {
	e.addBinding(captureBindingInfo{Name: name})
}

func (e *udfCaptureEnv) addBinding(binding captureBindingInfo) {
	if len(e.scopes) == 0 {
		return
	}

	e.scopes[len(e.scopes)-1][binding.Name] = binding
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
