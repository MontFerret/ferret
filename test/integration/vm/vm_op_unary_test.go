package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUnaryOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN !TRUE", false),
		S("RETURN NOT TRUE", false),
		S("RETURN !FALSE", true),
		S("RETURN -1", -1),
		S("RETURN -1.1", -1.1),
		S("RETURN +1", 1),
		S("RETURN +1.1", 1.1),
		S("LET v = 1 RETURN -v", -1),
		S("LET v = 1.1 RETURN -v", -1.1),
		S("LET v = -1 RETURN -v", 1),
		S("LET v = -1.1 RETURN -v", 1.1),
		S("LET v = -1 RETURN +v", -1),
		S("LET v = -1.1 RETURN +v", -1.1),
	})

	t.Run("RETURN { enabled: !val }", func(t *testing.T) {
		c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

		p1, err := c.Compile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !val }
		`))
		if err != nil {
			t.Fatalf("compile failed for single negation: %v", err)
		}

		vm1, err := vm.New(p1)
		if err != nil {
			t.Fatalf("vm init failed for single negation: %v", err)
		}

		r1, err := vm1.Run(context.Background(), nil)
		if err != nil {
			t.Fatalf("run failed for single negation: %v", err)
		}

		out1, err := encodingjson.Default.Encode(r1.Root())
		if closeErr := r1.Close(); closeErr != nil {
			t.Fatalf("close failed for single negation: %v", closeErr)
		}
		if err != nil {
			t.Fatalf("encode failed for single negation: %v", err)
		}

		if string(out1) != `{"enabled":true}` {
			t.Fatalf("unexpected single negation output: got %s", string(out1))
		}

		p2, err := c.Compile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !!val }
		`))
		if err != nil {
			t.Fatalf("compile failed for double negation: %v", err)
		}

		vm2, err := vm.New(p2)
		if err != nil {
			t.Fatalf("vm init failed for double negation: %v", err)
		}

		r2, err := vm2.Run(context.Background(), nil)
		if err != nil {
			t.Fatalf("run failed for double negation: %v", err)
		}

		out2, err := encodingjson.Default.Encode(r2.Root())
		if closeErr := r2.Close(); closeErr != nil {
			t.Fatalf("close failed for double negation: %v", closeErr)
		}
		if err != nil {
			t.Fatalf("encode failed for double negation: %v", err)
		}
		if string(out2) != `{"enabled":false}` {
			t.Fatalf("unexpected double negation output: got %s", string(out2))
		}
	})
}
