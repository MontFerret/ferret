// Package runtime provides v1-compatible runtime types for the Ferret compatibility layer.
// It mirrors the github.com/MontFerret/ferret/pkg/runtime package from Ferret v1.
package runtime

import (
	"context"
	"io"

	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Options holds the runtime execution options translated from v1-style Option functions.
type Options struct {
	logWriter io.Writer
	params    map[string]interface{}
	logFields map[string]interface{}
	logLevel  runtime.LogLevel
}

// Option is a functional option for configuring a program execution.
type Option func(*Options)

// WithParam sets a single named parameter for the query execution.
// The value is converted to a runtime.Value via runtime.ValueOf.
func WithParam(name string, value interface{}) Option {
	return func(o *Options) {
		if o.params == nil {
			o.params = make(map[string]interface{})
		}

		o.params[name] = value
	}
}

// WithParams sets multiple named parameters for the query execution.
func WithParams(params map[string]interface{}) Option {
	return func(o *Options) {
		if o.params == nil {
			o.params = make(map[string]interface{})
		}

		for k, v := range params {
			o.params[k] = v
		}
	}
}

// WithLog sets the writer for log output.
func WithLog(writer io.Writer) Option {
	return func(o *Options) {
		o.logWriter = writer
	}
}

// WithLogLevel sets the log level.
func WithLogLevel(lvl runtime.LogLevel) Option {
	return func(o *Options) {
		o.logLevel = lvl
	}
}

// WithLogFields sets additional structured log fields.
func WithLogFields(fields map[string]interface{}) Option {
	return func(o *Options) {
		o.logFields = fields
	}
}

// newOptions applies the provided setters to a default Options struct.
func newOptions(setters []Option) *Options {
	o := &Options{
		logLevel: runtime.ErrorLevel,
	}

	for _, s := range setters {
		if s != nil {
			s(o)
		}
	}

	return o
}

// ToSessionOptions converts a slice of compat Option setters into v2 SessionOptions.
// It is exported for use by the compat root package.
func ToSessionOptions(setters []Option) []ferret.SessionOption {
	return toSessionOptions(newOptions(setters))
}

// toSessionOptions converts compat Options into v2 SessionOptions.
func toSessionOptions(o *Options) []ferret.SessionOption {
	var opts []ferret.SessionOption

	if len(o.params) > 0 {
		params := make(runtime.Params, len(o.params))

		for k, v := range o.params {
			parsed, err := runtime.ValueOf(v)
			if err != nil {
				parsed = runtime.None
			}

			params[k] = parsed
		}

		opts = append(opts, ferret.WithSessionParams(params))
	}

	var envOpts []vm.EnvironmentOption

	if o.logWriter != nil {
		envOpts = append(envOpts, vm.WithLog(o.logWriter))
	}

	if o.logLevel != runtime.ErrorLevel {
		envOpts = append(envOpts, vm.WithLogLevel(o.logLevel))
	}

	if len(o.logFields) > 0 {
		envOpts = append(envOpts, vm.WithLogFields(o.logFields))
	}

	if len(envOpts) > 0 {
		opts = append(opts, ferret.WithEnvironmentOptions(envOpts...))
	}

	return opts
}

// Program is the v1-compatible compiled query program.
// It wraps a v2 Plan and exposes the v1 Program API.
type Program struct {
	plan   *ferret.Plan
	source string
	params []string
}

// NewProgram creates a compat Program from a v2 Plan.
// source is the original query string; params are the declared parameter names.
func NewProgram(plan *ferret.Plan, source string) *Program {
	return &Program{
		plan:   plan,
		source: source,
		params: plan.Params(),
	}
}

// Source returns the original query source text.
func (p *Program) Source() string {
	return p.source
}

// Params returns the list of parameter names declared in the query.
func (p *Program) Params() []string {
	return p.params
}

// Run executes the compiled program and returns the JSON-encoded result.
func (p *Program) Run(ctx context.Context, setters ...Option) ([]byte, error) {
	o := newOptions(setters)
	sessionOpts := toSessionOptions(o)

	session, err := p.plan.NewSession(ctx, sessionOpts...)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	out, err := session.Run(ctx)
	if err != nil {
		return nil, err
	}

	return out.Content, nil
}

// MustRun executes the program and panics on error.
func (p *Program) MustRun(ctx context.Context, setters ...Option) []byte {
	out, err := p.Run(ctx, setters...)
	if err != nil {
		panic(err)
	}

	return out
}

// compileFromSource is a helper used by compat packages to compile a raw query string
// using the provided engine and return a wrapped Program.
func CompileFromSource(ctx context.Context, eng *ferret.Engine, query string) (*Program, error) {
	src := file.NewAnonymousSource(query)

	plan, err := eng.Compile(ctx, src)
	if err != nil {
		return nil, err
	}

	return NewProgram(plan, query), nil
}
