package sdk_test

import (
	"context"
	"testing"

	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
	sdktest "github.com/MontFerret/ferret/v2/test/spec/sdk"
)

func TestHarnessRunsExternalModuleAndCleansUp(t *testing.T) {
	closed := false

	t.Run("external module", func(t *testing.T) {
		mod := sdk.NewModule("synthetic", func(bootstrap module.Bootstrap) error {
			bootstrap.Hooks().Engine().OnClose(func() error {
				closed = true
				return nil
			})

			return sdk.RegisterFunctions(
				bootstrap.Host().Library().Namespace("SYNTHETIC"),
				sdk.Func("ADD", sdk.Bind2(func(
					_ context.Context,
					left runtime.Int,
					right runtime.Int,
				) (runtime.Value, error) {
					return left + right, nil
				})),
			)
		})

		harness := sdktest.New(t, ferret.WithModules(mod))
		if harness.Engine() == nil {
			t.Fatal("expected harness engine")
		}

		output, err := harness.Run(t.Context(), "RETURN SYNTHETIC::ADD(2, 3)")
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		if string(output.Content) != "5" {
			t.Fatalf("got output %s", output.Content)
		}
	})

	if !closed {
		t.Fatal("expected testing cleanup to close the engine")
	}
}
