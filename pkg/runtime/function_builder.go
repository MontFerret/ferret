package runtime

type functionBuilder struct {
	functions Functions
}

func NewFunctionsBuilder() FunctionsBuilder {
	return &functionBuilder{
		functions: NewFunctions(),
	}
}

func (f *functionBuilder) Set(name string, fn Function) FunctionsBuilder {
	f.functions.F().Set(name, fn)

	return f
}

func (f *functionBuilder) Set0(name string, fn Function0) FunctionsBuilder {
	f.functions.F0().Set(name, fn)

	return f
}

func (f *functionBuilder) Set1(name string, fn Function1) FunctionsBuilder {
	f.functions.F1().Set(name, fn)

	return f
}

func (f *functionBuilder) Set2(name string, fn Function2) FunctionsBuilder {
	f.functions.F2().Set(name, fn)

	return f
}

func (f *functionBuilder) Set3(name string, fn Function3) FunctionsBuilder {
	f.functions.F3().Set(name, fn)

	return f
}

func (f *functionBuilder) Set4(name string, fn Function4) FunctionsBuilder {
	f.functions.F4().Set(name, fn)

	return f
}

func (f *functionBuilder) Build() Functions {
	return f.functions
}
