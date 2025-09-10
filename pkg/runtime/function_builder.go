package runtime

type (
	FunctionsBuilder interface {
		Has(name string) bool
		Set(name string, fn Function) FunctionsBuilder
		Set0(name string, fn Function0) FunctionsBuilder
		Set1(name string, fn Function1) FunctionsBuilder
		Set2(name string, fn Function2) FunctionsBuilder
		Set3(name string, fn Function3) FunctionsBuilder
		Set4(name string, fn Function4) FunctionsBuilder
		SetFrom(other Functions) FunctionsBuilder
		Unset(name string) FunctionsBuilder
		Unset0(name string) FunctionsBuilder
		Unset1(name string) FunctionsBuilder
		Unset2(name string) FunctionsBuilder
		Unset3(name string) FunctionsBuilder
		Unset4(name string) FunctionsBuilder
		UnsetFrom(other Functions) FunctionsBuilder
		Build() Functions
	}

	defaultFunctionBuilder struct {
		namespace string
		f         map[string]Function
		f0        map[string]Function0
		f1        map[string]Function1
		f2        map[string]Function2
		f3        map[string]Function3
		f4        map[string]Function4
	}
)

func NewFunctionsBuilder() FunctionsBuilder {
	return &defaultFunctionBuilder{}
}

func newRootFunctionsBuilder() *defaultFunctionBuilder {
	return newNamespaceFunctionsBuilder("")
}

func newNamespaceFunctionsBuilder(namespace string) *defaultFunctionBuilder {
	return &defaultFunctionBuilder{
		namespace: namespace,
		f:         make(map[string]Function),
		f0:        make(map[string]Function0),
		f1:        make(map[string]Function1),
		f2:        make(map[string]Function2),
		f3:        make(map[string]Function3),
		f4:        make(map[string]Function4),
	}
}

func newFunctionsBuilderInternalFrom(namespace string, other *defaultFunctionBuilder) *defaultFunctionBuilder {
	b := newNamespaceFunctionsBuilder(namespace)
	b.f = other.f
	b.f0 = other.f0
	b.f1 = other.f1
	b.f2 = other.f2
	b.f3 = other.f3
	b.f4 = other.f4

	return b
}

func (b *defaultFunctionBuilder) Has(name string) bool {
	fname := makeFunctionName(b.namespace, name)

	if b.f != nil {
		if _, ok := b.f[fname]; ok {
			return true
		}
	}

	if b.f0 != nil {
		if _, ok := b.f0[fname]; ok {
			return true
		}
	}

	if b.f1 != nil {
		if _, ok := b.f1[fname]; ok {
			return true
		}
	}

	if b.f2 != nil {
		if _, ok := b.f2[fname]; ok {
			return true
		}
	}

	if b.f3 != nil {
		if _, ok := b.f3[fname]; ok {
			return true
		}
	}

	if b.f4 != nil {
		if _, ok := b.f4[fname]; ok {
			return true
		}
	}

	return false
}

func (b *defaultFunctionBuilder) Set(name string, fn Function) FunctionsBuilder {
	if b.f == nil {
		b.f = make(map[string]Function)
	}

	b.f[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) Set0(name string, fn Function0) FunctionsBuilder {
	if b.f0 == nil {
		b.f0 = make(map[string]Function0)
	}

	b.f0[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) Set1(name string, fn Function1) FunctionsBuilder {
	if b.f1 == nil {
		b.f1 = make(map[string]Function1)
	}

	b.f1[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) Set2(name string, fn Function2) FunctionsBuilder {
	if b.f2 == nil {
		b.f2 = make(map[string]Function2)
	}

	b.f2[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) Set3(name string, fn Function3) FunctionsBuilder {
	if b.f3 == nil {
		b.f3 = make(map[string]Function3)
	}

	b.f3[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) Set4(name string, fn Function4) FunctionsBuilder {
	if b.f4 == nil {
		b.f4 = make(map[string]Function4)
	}

	b.f4[makeFunctionName(b.namespace, name)] = fn

	return b
}

func (b *defaultFunctionBuilder) SetFrom(other Functions) FunctionsBuilder {
	_ = other.F().ForEach(func(fn Function, name string) error {
		b.Set(name, fn)

		return nil
	})

	_ = other.F0().ForEach(func(fn Function0, name string) error {
		b.Set0(name, fn)

		return nil
	})

	_ = other.F1().ForEach(func(fn Function1, name string) error {
		b.Set1(name, fn)

		return nil
	})

	_ = other.F2().ForEach(func(fn Function2, name string) error {
		b.Set2(name, fn)

		return nil
	})

	_ = other.F3().ForEach(func(fn Function3, name string) error {
		b.Set3(name, fn)

		return nil
	})

	_ = other.F4().ForEach(func(fn Function4, name string) error {
		b.Set4(name, fn)

		return nil
	})

	return b
}

func (b *defaultFunctionBuilder) Unset(name string) FunctionsBuilder {
	if b.f != nil {
		delete(b.f, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) Unset0(name string) FunctionsBuilder {
	if b.f0 != nil {
		delete(b.f0, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) Unset1(name string) FunctionsBuilder {
	if b.f1 != nil {
		delete(b.f1, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) Unset2(name string) FunctionsBuilder {
	if b.f2 != nil {
		delete(b.f2, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) Unset3(name string) FunctionsBuilder {
	if b.f3 != nil {
		delete(b.f3, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) Unset4(name string) FunctionsBuilder {
	if b.f4 != nil {
		delete(b.f4, makeFunctionName(b.namespace, name))
	}

	return b
}

func (b *defaultFunctionBuilder) UnsetFrom(other Functions) FunctionsBuilder {
	_ = other.F().ForEach(func(_ Function, name string) error {
		b.Unset(name)

		return nil
	})

	_ = other.F0().ForEach(func(_ Function0, name string) error {
		b.Unset0(name)

		return nil
	})

	_ = other.F1().ForEach(func(_ Function1, name string) error {
		b.Unset1(name)

		return nil
	})

	_ = other.F2().ForEach(func(_ Function2, name string) error {
		b.Unset2(name)

		return nil
	})

	_ = other.F3().ForEach(func(_ Function3, name string) error {
		b.Unset3(name)

		return nil
	})

	_ = other.F4().ForEach(func(_ Function4, name string) error {
		b.Unset4(name)

		return nil
	})

	return b
}

func (b *defaultFunctionBuilder) Build() Functions {
	registry := &functionRegistry{}

	if len(b.f) > 0 {
		registry.f = NewFunctionCollectionFromMap(b.f)
	}

	if len(b.f0) > 0 {
		registry.f0 = NewFunctionCollectionFromMap(b.f0)
	}

	if len(b.f1) > 0 {
		registry.f1 = NewFunctionCollectionFromMap(b.f1)
	}

	if len(b.f2) > 0 {
		registry.f2 = NewFunctionCollectionFromMap(b.f2)
	}

	if len(b.f3) > 0 {
		registry.f3 = NewFunctionCollectionFromMap(b.f3)
	}

	if len(b.f4) > 0 {
		registry.f4 = NewFunctionCollectionFromMap(b.f4)
	}

	return registry
}
