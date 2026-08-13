package analyzer

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRegisteredSignaturesRejectsNilFunction(t *testing.T) {
	builder := runtime.NewFunctionsBuilder()
	var fn runtime.Function0
	builder.A0().Add("nil", fn)

	functions, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}

	_, err = registeredSignatures(functions)
	if err == nil || !strings.Contains(err.Error(), "nil/0") || !strings.Contains(err.Error(), "nil") {
		t.Fatalf("error = %v, want source-aware nil registration error", err)
	}
}

func TestValidateRegisteredNameRequiresCanonicalQualifiedName(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"lower", "db::postgres::query"} {
		if err := validateRegisteredName(name); err != nil {
			t.Errorf("validateRegisteredName(%q): %v", name, err)
		}
	}

	for _, name := range []string{"", "db::", "UPPER", "DB::postgres::query", "db::POSTGRES::query", "db::postgres::QuErY"} {
		if err := validateRegisteredName(name); err == nil {
			t.Errorf("validateRegisteredName(%q) unexpectedly succeeded", name)
		}
	}
}
