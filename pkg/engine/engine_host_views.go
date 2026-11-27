package engine

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	codecResolverView struct {
		registry *encoding.Registry
	}

	paramsReaderView struct {
		params runtime.Params
	}

	functionRegistryView struct {
		functions *runtime.Functions
	}
)

// Encoding returns the reusable read-only view of the engine's finalized codec
// registrations. Resolved codec instances are shared with the engine.
func (e *Engine) Encoding() encoding.CodecResolver {
	return e.codecView
}

// Params returns the reusable read-only view of the engine's default parameter
// bindings. The retrieved values are shared and may be mutable.
func (e *Engine) Params() runtime.ParamsReader {
	return e.paramsView
}

// Functions returns the reusable read-only view of the engine's finalized
// function registry.
func (e *Engine) Functions() runtime.FunctionRegistry {
	return e.functionsView
}

func (v *codecResolverView) Codec(contentType string) (encoding.Codec, error) {
	return v.registry.Codec(contentType)
}

func (v *codecResolverView) Encoder(contentType string) (encoding.Encoder, error) {
	return v.registry.Encoder(contentType)
}

func (v *codecResolverView) Decoder(contentType string) (encoding.Decoder, error) {
	return v.registry.Decoder(contentType)
}

func (v *paramsReaderView) Get(name string) (runtime.Value, bool) {
	return v.params.Get(name)
}

func (v *paramsReaderView) MustGet(name string) runtime.Value {
	return v.params.MustGet(name)
}

func (v *paramsReaderView) GetOr(name string, fallback runtime.Value) runtime.Value {
	return v.params.GetOr(name, fallback)
}

func (v *paramsReaderView) Has(name string) bool {
	return v.params.Has(name)
}

func (v *paramsReaderView) Size() int {
	return v.params.Size()
}

func (v *paramsReaderView) Names() []string {
	return v.params.Names()
}

func (v *functionRegistryView) Hash() uint64 {
	return v.functions.Hash()
}

func (v *functionRegistryView) Size() int {
	return v.functions.Size()
}

func (v *functionRegistryView) Has(name string) bool {
	return v.functions.Has(name)
}

func (v *functionRegistryView) List() []string {
	return v.functions.List()
}

func (v *functionRegistryView) A0() runtime.FunctionCollection[runtime.Function0] {
	return v.functions.A0()
}

func (v *functionRegistryView) A1() runtime.FunctionCollection[runtime.Function1] {
	return v.functions.A1()
}

func (v *functionRegistryView) A2() runtime.FunctionCollection[runtime.Function2] {
	return v.functions.A2()
}

func (v *functionRegistryView) A3() runtime.FunctionCollection[runtime.Function3] {
	return v.functions.A3()
}

func (v *functionRegistryView) A4() runtime.FunctionCollection[runtime.Function4] {
	return v.functions.A4()
}

func (v *functionRegistryView) Var() runtime.FunctionCollection[runtime.Function] {
	return v.functions.Var()
}
