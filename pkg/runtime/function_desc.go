package runtime

type (
	FunctionOption func(*FunctionDescriptor)

	FunctionDescriptor struct {
		Name        string
		Summary     string
		Description string
		Params      []FunctionParamDescriptor
		Return      FunctionReturnDescriptor
		Variadic    bool
	}

	FunctionParamDescriptor struct {
		Name     string
		Type     []Type
		Summary  string
		Optional bool
	}

	FunctionReturnDescriptor struct {
		Type    []Type
		Summary string
	}
)

func WithSummary(value string) FunctionOption {
	return func(d *FunctionDescriptor) {
		d.Summary = value
	}
}

func WithDescription(value string) FunctionOption {
	return func(d *FunctionDescriptor) {
		d.Description = value
	}
}

func WithParam(name string, typ ...Type) FunctionOption {
	return func(d *FunctionDescriptor) {
		d.Params = append(d.Params, FunctionParamDescriptor{
			Name: name,
			Type: typ,
		})
	}
}

func WithOptionalParam(name string, typ ...Type) FunctionOption {
	return func(d *FunctionDescriptor) {
		d.Params = append(d.Params, FunctionParamDescriptor{
			Name:     name,
			Type:     typ,
			Optional: true,
		})
	}
}

func WithReturn(typ ...Type) FunctionOption {
	return func(d *FunctionDescriptor) {
		d.Return.Type = typ
	}
}
