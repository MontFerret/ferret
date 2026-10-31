package format

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func TestRunnerInvokesDeclaredAssertions(t *testing.T) {
	f, err := formatter.New()
	if err != nil {
		t.Fatal(err)
	}

	var outputs, errors []any
	recordOutput := func(_ *testing.T, actual any, _ ...any) {
		outputs = append(outputs, actual)
	}
	recordError := func(_ *testing.T, actual any, _ ...any) {
		errors = append(errors, actual)
	}

	runner := Runner{Name: "Formatter", Formatter: f}
	runner.Run(t, []Spec{
		NewSpec("RETURN   1").Expect().Output(recordOutput).Expect().OutputError(recordError),
		NewSpec("RETURN [").Expect().Output(assert.ShouldNotReach).Expect().OutputError(recordError),
	})

	assert.ShouldEqual(t, outputs, []any{"return 1"})
	if len(errors) != 2 {
		t.Fatalf("error assertion called %d times, want 2", len(errors))
	}

	assert.ShouldBeNil(t, errors[0])
	assert.ShouldBeError(t, errors[1])
}
