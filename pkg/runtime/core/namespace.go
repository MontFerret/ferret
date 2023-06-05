package core

type Namespace interface {
	Namespace(name string) Namespace
	RegisterFunction(name string, fun Function) error
	RegisterFunctions(funs Functions) error
	RegisteredFunctions() []string
	RemoveFunction(name string)
}
