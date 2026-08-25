package ferret_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MontFerret/api"
	apisource "github.com/MontFerret/api/source"
	ferret "github.com/MontFerret/ferret/v2"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
)

var (
	_ ferret.Source  = apisource.File{}
	_ apisource.File = ferret.Source{}

	_ ferret.Output = api.Output{}
	_ api.Output    = ferret.Output{}

	_ ferret.PlanOption  = api.PlanOption(nil)
	_ api.PlanOption     = ferret.PlanOption(nil)
	_ ferret.PlanOptions = api.PlanOptions(nil)
	_ api.PlanOptions    = ferret.PlanOptions(nil)

	_ ferret.SessionOption  = api.SessionOption(nil)
	_ api.SessionOption     = ferret.SessionOption(nil)
	_ ferret.SessionOptions = api.SessionOptions(nil)
	_ api.SessionOptions    = ferret.SessionOptions(nil)

	_ ferret.OptimizationLevel = api.OptimizationNone
	_ api.OptimizationLevel    = ferret.OptimizationNone
	_ api.OptimizationLevel    = ferret.OptimizationBasic
	_ api.OptimizationLevel    = ferret.OptimizationFull
	_ api.OptimizationLevel    = ferret.OptimizationAggressive
)

type foreignPlanOptions struct {
	level api.OptimizationLevel
}

func (o *foreignPlanOptions) SetOptimizationLevel(level api.OptimizationLevel) error {
	o.level = level

	return nil
}

func TestSourceConstructorsForwardToUniversalAPI(t *testing.T) {
	t.Parallel()

	if got, want := ferret.NewSource("main.fql", "RETURN 1"), apisource.New("main.fql", "RETURN 1"); got != want {
		t.Fatalf("named source = %#v, want %#v", got, want)
	}

	if got, want := ferret.NewAnonymousSource("RETURN 1"), apisource.NewAnonymous("RETURN 1"); got != want {
		t.Fatalf("anonymous source = %#v, want %#v", got, want)
	}
}

func TestFacadeAliasesHaveExactUniversalIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		got  reflect.Type
		want reflect.Type
		name string
	}{
		{name: "source", got: reflect.TypeFor[ferret.Source](), want: reflect.TypeFor[apisource.File]()},
		{name: "output", got: reflect.TypeFor[ferret.Output](), want: reflect.TypeFor[api.Output]()},
		{name: "plan option", got: reflect.TypeFor[ferret.PlanOption](), want: reflect.TypeFor[api.PlanOption]()},
		{name: "plan options", got: reflect.TypeFor[ferret.PlanOptions](), want: reflect.TypeFor[api.PlanOptions]()},
		{name: "session option", got: reflect.TypeFor[ferret.SessionOption](), want: reflect.TypeFor[api.SessionOption]()},
		{name: "session options", got: reflect.TypeFor[ferret.SessionOptions](), want: reflect.TypeFor[api.SessionOptions]()},
		{name: "optimization level", got: reflect.TypeFor[ferret.OptimizationLevel](), want: reflect.TypeFor[api.OptimizationLevel]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.got != tt.want {
				t.Fatalf("type = %v, want exact identity with %v", tt.got, tt.want)
			}
		})
	}
}

func TestOptimizationOptionForwardsToForeignTarget(t *testing.T) {
	t.Parallel()

	target := &foreignPlanOptions{}
	if err := ferret.WithOptimizationLevel(ferret.OptimizationAggressive)(target); err != nil {
		t.Fatalf("apply optimization option: %v", err)
	}

	if target.level != api.OptimizationAggressive {
		t.Fatalf("optimization level = %d, want %d", target.level, api.OptimizationAggressive)
	}
}

func TestNativeFacadeOptionsExecuteEndToEnd(t *testing.T) {
	t.Parallel()

	engine, err := ferret.New(ferret.WithEngineParam("value", "engine"))
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	plan, err := engine.Compile(
		context.Background(),
		ferret.NewSource("facade.fql", "RETURN [@value, @other]"),
		ferret.WithOptimizationLevel(ferret.OptimizationFull),
	)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewSession(
		context.Background(),
		ferret.WithParams(map[string]any{
			"value": "map",
			"other": 2,
		}),
		ferret.WithParam("value", "session"),
		ferret.WithOutputContentType(encodingjson.ContentType),
	)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	output, err := session.Run(context.Background())
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	if output.ContentType != encodingjson.ContentType {
		t.Fatalf("content type = %q, want %q", output.ContentType, encodingjson.ContentType)
	}

	if got, want := string(output.Content), `["session",2]`; got != want {
		t.Fatalf("content = %s, want %s", got, want)
	}
}
