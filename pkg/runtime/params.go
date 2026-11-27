package runtime

import (
	"fmt"
	"slices"
)

type (
	// ParamsReader reads parameter bindings. Retrieved values may be mutable.
	ParamsReader interface {
		Get(name string) (Value, bool)
		MustGet(name string) Value
		GetOr(name string, fallback Value) Value
		Has(name string) bool
		Size() int
		Names() []string
	}

	Params map[string]Value
)

var _ ParamsReader = Params(nil)

// NewParams creates and returns an empty Params map.
func NewParams() Params {
	return make(Params)
}

// NewParamsFrom constructs a Params object from a map of string keys to any values, converting each value to a Value.
func NewParamsFrom(values map[string]any) (Params, error) {
	params := make(Params, len(values))

	for k, v := range values {
		val, err := ValueOf(v)

		if err != nil {
			return nil, fmt.Errorf("param %q: %w", k, err)
		}

		params[k] = val
	}

	return params, nil
}

func (p Params) Get(name string) (Value, bool) {
	value, exists := p[name]

	if !exists {
		return None, false
	}

	if value == nil {
		return None, true
	}

	return value, true
}

func (p Params) MustGet(name string) Value {
	value, exists := p.Get(name)

	if !exists {
		panic(Errorf(ErrNotFound, "param %q", name))
	}

	return value
}

func (p Params) GetOr(name string, fallback Value) Value {
	value, exists := p.Get(name)

	if exists {
		return value
	}

	if fallback == nil {
		return None
	}

	return fallback
}

func (p Params) Has(name string) bool {
	_, exists := p[name]

	return exists
}

// Size returns the number of bindings.
func (p Params) Size() int {
	return len(p)
}

// Names returns the parameter names in lexicographic order as a new slice.
func (p Params) Names() []string {
	names := make([]string, 0, len(p))
	for name := range p {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}

func (p Params) SetValue(name string, value Value) Params {
	p[name] = value

	return p
}

func (p Params) SetAllValues(values map[string]Value) Params {
	for k, v := range values {
		p[k] = v
	}

	return p
}

func (p Params) Set(name string, value any) error {
	v, err := ValueOf(value)

	if err != nil {
		return fmt.Errorf("param %q: %w", name, err)
	}

	p[name] = v

	return nil
}

func (p Params) SetAll(values map[string]any) error {
	for k, v := range values {
		if err := p.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (p Params) Merge(other map[string]any) (Params, error) {
	if p == nil {
		return NewParamsFrom(other)
	}

	if other == nil {
		return p, nil
	}

	for k, v := range other {
		if err := p.Set(k, v); err != nil {
			return p, err
		}
	}

	return p, nil
}

func (p Params) MergeParams(other Params) Params {
	if p == nil {
		return NewParams().MergeParams(other)
	}

	if other == nil {
		return p
	}

	for k, v := range other {
		p[k] = v
	}

	return p
}

func (p Params) MustSet(name string, value any) Params {
	if err := p.Set(name, value); err != nil {
		panic(err)
	}

	return p
}

func (p Params) Delete(name string) {
	delete(p, name)
}

func (p Params) Clone() Params {
	if p == nil {
		return nil
	}

	out := make(Params, len(p))

	for k, v := range p {
		out[k] = v
	}

	return out
}
