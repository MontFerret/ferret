package universal

import (
	"errors"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type sessionOptions struct {
	native []engine.SessionOption
}

var _ api.SessionOptions = (*sessionOptions)(nil)

func (o *sessionOptions) SetParam(name string, value any) error {
	o.native = append(o.native, engine.WithSessionParam(name, value))

	return nil
}

func (o *sessionOptions) SetParams(params map[string]any) error {
	o.native = append(o.native, engine.WithSessionParams(params))

	return nil
}

func (o *sessionOptions) SetOutputContentType(value string) error {
	o.native = append(o.native, engine.WithOutputContentType(value))

	return nil
}

func (o *sessionOptions) SetFSRoot(value string) error {
	o.native = append(o.native, engine.WithSessionFSRoot(value))

	return nil
}

func newSessionOptions(setters []api.SessionOption) (*sessionOptions, error) {
	opts := &sessionOptions{}
	var failures []error
	for _, setter := range setters {
		if setter == nil {
			continue
		}

		if err := setter(opts); err != nil {
			failures = append(failures, err)
		}
	}

	return opts, errors.Join(failures...)
}
