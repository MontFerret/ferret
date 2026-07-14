package sdk_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

func TestNewModule(t *testing.T) {
	t.Run("registers through callback", func(t *testing.T) {
		called := false
		mod := sdk.NewModule("example", func(_ module.Bootstrap) error {
			called = true
			return nil
		})

		if mod.Name() != "example" {
			t.Fatalf("unexpected name %q", mod.Name())
		}

		if err := mod.Register(nil); err != nil {
			t.Fatalf("unexpected registration error: %v", err)
		}

		if !called {
			t.Fatal("expected registration callback to be called")
		}
	})

	t.Run("defers invalid constructor input", func(t *testing.T) {
		if err := sdk.NewModule("", func(_ module.Bootstrap) error { return nil }).Register(nil); err == nil {
			t.Fatal("expected empty name error")
		}

		if err := sdk.NewModule("example", nil).Register(nil); err == nil {
			t.Fatal("expected nil callback error")
		}
	})

	t.Run("adds module context to callback errors", func(t *testing.T) {
		sentinel := errors.New("failed")
		mod := sdk.NewModule("example", func(_ module.Bootstrap) error {
			return sentinel
		})

		err := mod.Register(nil)
		if !errors.Is(err, sentinel) {
			t.Fatalf("expected wrapped sentinel, got %v", err)
		}

		if !strings.Contains(err.Error(), `module "example": register`) {
			t.Fatalf("expected module context, got %v", err)
		}

		engine, err := ferret.New(ferret.WithModules(mod))
		if engine != nil || !errors.Is(err, sentinel) {
			t.Fatalf("expected engine registration failure, got engine=%v err=%v", engine, err)
		}
	})
}

func TestRegisterFunctions(t *testing.T) {
	library := runtime.NewLibrary()
	ns := library.Namespace("TEST")

	err := sdk.RegisterFunctions(ns,
		sdk.Func("ZERO", runtime.Function0(func(context.Context) (runtime.Value, error) {
			return runtime.None, nil
		})),
		sdk.Func("ONE", runtime.Function1(func(context.Context, runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})),
		sdk.Func("TWO", runtime.Function2(func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})),
		sdk.Func("THREE", runtime.Function3(func(context.Context, runtime.Value, runtime.Value, runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})),
		sdk.Func("FOUR", runtime.Function4(func(context.Context, runtime.Value, runtime.Value, runtime.Value, runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})),
		sdk.Func("VARIABLE", runtime.Function(func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})),
	)
	if err != nil {
		t.Fatalf("register functions: %v", err)
	}

	functions, err := library.Build()
	if err != nil {
		t.Fatalf("build library: %v", err)
	}

	for _, name := range []string{"ZERO", "ONE", "TWO", "THREE", "FOUR", "VARIABLE"} {
		if !functions.Has("TEST::" + name) {
			t.Errorf("expected %s to be registered", name)
		}
	}
}

func TestRegisterFunctionsIsAtomic(t *testing.T) {
	library := runtime.NewLibrary()
	ns := library.Namespace("TEST")
	valid := runtime.Function0(func(context.Context) (runtime.Value, error) {
		return runtime.None, nil
	})

	err := sdk.RegisterFunctions(ns,
		sdk.Func("VALID", valid),
		sdk.Func("", valid),
	)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if ns.Function().Has("VALID") {
		t.Fatal("valid definition was registered after atomic validation failed")
	}

	err = sdk.RegisterFunctions(ns,
		sdk.Func("DUPLICATE", valid),
		sdk.Func("DUPLICATE", valid),
	)
	if err == nil {
		t.Fatal("expected duplicate error")
	}

	if ns.Function().Has("DUPLICATE") {
		t.Fatal("duplicate definition was partially registered")
	}

	if err := sdk.RegisterFunctions(ns, sdk.Func("EXISTING", valid)); err != nil {
		t.Fatalf("register existing function: %v", err)
	}
	err = sdk.RegisterFunctions(ns,
		sdk.Func("NEW", valid),
		sdk.Func("EXISTING", valid),
	)
	if err == nil {
		t.Fatal("expected existing definition error")
	}
	if ns.Function().Has("NEW") {
		t.Fatal("new definition was registered before existing-name validation failed")
	}

	var nilFunction runtime.Function0
	if err := sdk.RegisterFunctions(ns, sdk.Func("NIL", nilFunction)); err == nil {
		t.Fatal("expected nil function error")
	}
}

func TestTypedBinders(t *testing.T) {
	ctx := t.Context()

	zero := sdk.Bind0(func(context.Context) (runtime.Value, error) {
		return runtime.NewString("zero"), nil
	})
	one := sdk.Bind1(func(_ context.Context, value runtime.String) (runtime.Value, error) {
		return value, nil
	})
	two := sdk.Bind2(func(_ context.Context, left runtime.String, right runtime.Int) (runtime.Value, error) {
		return runtime.NewString(left.String() + right.String()), nil
	})
	three := sdk.Bind3(func(_ context.Context, first, second, third runtime.Int) (runtime.Value, error) {
		return first + second + third, nil
	})
	four := sdk.Bind4(func(_ context.Context, first, second, third, fourth runtime.Int) (runtime.Value, error) {
		return first + second + third + fourth, nil
	})

	value, err := zero(ctx)
	assertBoundResult(t, value, err, "zero")
	value, err = one(ctx, runtime.NewString("one"))
	assertBoundResult(t, value, err, "one")
	value, err = two(ctx, runtime.NewString("value"), runtime.NewInt(2))
	assertBoundResult(t, value, err, "value2")
	value, err = three(ctx, runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3))
	assertBoundResult(t, value, err, "6")
	value, err = four(ctx, runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3), runtime.NewInt(4))
	assertBoundResult(t, value, err, "10")
}

func TestTypedBinderErrorsAndNilResults(t *testing.T) {
	bound := sdk.Bind2(func(_ context.Context, left runtime.String, right runtime.Int) (runtime.Value, error) {
		return runtime.NewString(left.String() + right.String()), nil
	})

	_, err := bound(t.Context(), runtime.NewString("ok"), runtime.NewString("wrong"))
	position, ok, _ := runtime.InvalidArgumentDetails(err)
	if !ok || position != 1 {
		t.Fatalf("expected argument position 1, got position=%d ok=%v err=%v", position, ok, err)
	}

	sentinel := errors.New("operation failed")
	failing := sdk.Bind0(func(context.Context) (runtime.Value, error) {
		return runtime.NewString("ignored"), sentinel
	})
	value, err := failing(t.Context())
	if value != runtime.None || !errors.Is(err, sentinel) {
		t.Fatalf("expected None and sentinel, got value=%v err=%v", value, err)
	}

	nilResult := sdk.Bind0(func(context.Context) (runtime.Value, error) {
		return nil, nil
	})
	value, err = nilResult(t.Context())
	if err != nil || value != runtime.None {
		t.Fatalf("expected nil result to normalize to None, got value=%v err=%v", value, err)
	}

	typedNilResult := sdk.Bind0(func(context.Context) (*sdk.HostValue[int], error) {
		return nil, nil
	})
	value, err = typedNilResult(t.Context())
	if err != nil || value != runtime.None {
		t.Fatalf("expected typed nil result to normalize to None, got value=%v err=%v", value, err)
	}
}

func TestTypedBinderArgumentPositions(t *testing.T) {
	bound := sdk.Bind4(func(_ context.Context, first, second, third, fourth runtime.Int) (runtime.Value, error) {
		return first + second + third + fourth, nil
	})

	for expected := 0; expected < 4; expected++ {
		args := []runtime.Value{
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		}
		args[expected] = runtime.NewString("wrong")

		_, err := bound(t.Context(), args[0], args[1], args[2], args[3])
		position, ok, _ := runtime.InvalidArgumentDetails(err)
		if !ok || position != expected {
			t.Fatalf("expected argument position %d, got position=%d ok=%v err=%v", expected, position, ok, err)
		}
	}
}

func assertBoundResult(t *testing.T, value runtime.Value, err error, expected string) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}

	if value.String() != expected {
		t.Fatalf("got %q, want %q", value.String(), expected)
	}
}
