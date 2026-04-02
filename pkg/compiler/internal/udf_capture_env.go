package internal

import "github.com/antlr4-go/antlr/v4"

type (
	udfCaptureEnv struct {
		scopes []map[string]udfCaptureBinding
	}

	udfCaptureBinding struct {
		Decl    antlr.ParserRuleContext
		Name    string
		Mutable bool
	}
)

func (e *udfCaptureEnv) push() {
	e.scopes = append(e.scopes, make(map[string]udfCaptureBinding))
}

func (e *udfCaptureEnv) pop() {
	if len(e.scopes) > 0 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *udfCaptureEnv) add(name string) {
	e.addBinding(udfCaptureBinding{Name: name})
}

func (e *udfCaptureEnv) addBinding(binding udfCaptureBinding) {
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

func (e *udfCaptureEnv) resolveBinding(name string) (udfCaptureBinding, bool) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if binding, ok := e.scopes[i][name]; ok {
			return binding, true
		}
	}

	return udfCaptureBinding{}, false
}
