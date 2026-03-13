package vm_test

func RuntimeErrorCase(expression string, expected ExpectedRuntimeError, desc ...string) UseCase {
	return NewCase(expression, &expected, ShouldBeRuntimeError, desc...)
}

func SkipRuntimeErrorCase(expression string, expected ExpectedRuntimeError, desc ...string) UseCase {
	return Skip(RuntimeErrorCase(expression, expected, desc...))
}

func DebugRuntimeErrorCase(expression string, expected ExpectedRuntimeError, desc ...string) UseCase {
	return Debug(RuntimeErrorCase(expression, expected, desc...))
}
