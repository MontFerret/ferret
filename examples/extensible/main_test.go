package main

import (
	"testing"

	ferret "github.com/MontFerret/ferret/v2"
	sdktest "github.com/MontFerret/ferret/v2/test/spec/sdk"
)

func TestExampleModule(t *testing.T) {
	harness := sdktest.New(t, ferret.WithModules(newExampleModule()))

	output, err := harness.Run(
		t.Context(),
		`RETURN EXAMPLE::FORMAT(EXAMPLE::UPPER("ferret"), { prefix: "hello " })`,
	)
	if err != nil {
		t.Fatalf("run example module: %v", err)
	}
	if string(output.Content) != `"hello FERRET"` {
		t.Fatalf("unexpected output %s", output.Content)
	}
}

func TestExampleModuleRejectsUnknownOptions(t *testing.T) {
	harness := sdktest.New(t, ferret.WithModules(newExampleModule()))

	_, err := harness.Run(
		t.Context(),
		`RETURN EXAMPLE::FORMAT("ferret", { prefix: "hello ", typo: true })`,
	)
	if err == nil {
		t.Fatalf("expected strict option error, got %v", err)
	}
}
