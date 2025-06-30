package runtime

type Namespace interface {
	Namespace(name string) Namespace
	RegisterFunctions(funs Functions) error
	RegisteredFunctions() []string
	RemoveFunction(name string)
}
