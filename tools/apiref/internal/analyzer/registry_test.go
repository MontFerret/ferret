package analyzer

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRegisteredSignaturesRejectsNilFunction(t *testing.T) {
	builder := runtime.NewFunctionsBuilder()
	var fn runtime.Function0
	builder.A0().Add("NIL", fn)

	functions, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}

	_, err = registeredSignatures(functions)
	if err == nil || !strings.Contains(err.Error(), "NIL/0") || !strings.Contains(err.Error(), "nil") {
		t.Fatalf("error = %v, want source-aware nil registration error", err)
	}
}
