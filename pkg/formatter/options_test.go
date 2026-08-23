package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/formatter/internal"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func mustNewFormatter(t testing.TB, setters ...Option) *Formatter {
	t.Helper()

	formatterInstance, err := New(setters...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return formatterInstance
}

func TestFormatterOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		formatterInstance := mustNewFormatter(t)
		want := internal.DefaultConfig()

		if formatterInstance.config != want {
			t.Fatalf("config = %+v, want %+v", formatterInstance.config, want)
		}
		if formatterInstance.config.PrintWidth != 80 {
			t.Fatalf("print width = %d, want 80", formatterInstance.config.PrintWidth)
		}
		if formatterInstance.config.TabWidth != 4 {
			t.Fatalf("tab width = %d, want 4", formatterInstance.config.TabWidth)
		}
		if formatterInstance.config.SingleQuote {
			t.Fatal("single quote = true, want false")
		}
		if !formatterInstance.config.BracketSpacing {
			t.Fatal("bracket spacing = false, want true")
		}
		if formatterInstance.config.CaseMode != CaseModeLower {
			t.Fatalf("case mode = %d, want %d", formatterInstance.config.CaseMode, CaseModeLower)
		}
	})

	t.Run("valid widths", func(t *testing.T) {
		formatterInstance := mustNewFormatter(
			t,
			WithPrintWidth(120),
			WithTabWidth(2),
		)

		if formatterInstance.config.PrintWidth != 120 {
			t.Fatalf("print width = %d, want 120", formatterInstance.config.PrintWidth)
		}
		if formatterInstance.config.TabWidth != 2 {
			t.Fatalf("tab width = %d, want 2", formatterInstance.config.TabWidth)
		}

		got := formatForOptionTest(t, formatterInstance, "RETURN FOR value IN [1] { RETURN value }")
		want := "return for value in [1] {\n  return value\n}"
		if got != want {
			t.Fatalf("formatted output = %q, want %q", got, want)
		}
	})

	t.Run("boolean options", func(t *testing.T) {
		tests := []struct {
			option Option
			name   string
			input  string
			want   string
		}{
			{
				name:   "single quote enabled",
				input:  `RETURN "value"`,
				want:   "return 'value'",
				option: WithSingleQuote(true),
			},
			{
				name:   "single quote disabled",
				input:  `RETURN 'value'`,
				want:   `return "value"`,
				option: WithSingleQuote(false),
			},
			{
				name:   "bracket spacing enabled",
				input:  `RETURN {value:1}`,
				want:   `return { value: 1 }`,
				option: WithBracketSpacing(true),
			},
			{
				name:   "bracket spacing disabled",
				input:  `RETURN { value: 1 }`,
				want:   `return {value: 1}`,
				option: WithBracketSpacing(false),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				formatterInstance := mustNewFormatter(t, tt.option)
				got := formatForOptionTest(t, formatterInstance, tt.input)
				if got != tt.want {
					t.Fatalf("formatted output = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("case modes", func(t *testing.T) {
		tests := []struct {
			name string
			want string
			mode CaseMode
		}{
			{name: "ignore", mode: CaseModeIgnore, want: "RETURN TrUe"},
			{name: "upper", mode: CaseModeUpper, want: "RETURN TRUE"},
			{name: "lower", mode: CaseModeLower, want: "return true"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				formatterInstance := mustNewFormatter(t, WithCaseMode(tt.mode))
				if formatterInstance.config.CaseMode != tt.mode {
					t.Fatalf("case mode = %d, want %d", formatterInstance.config.CaseMode, tt.mode)
				}

				got := formatForOptionTest(t, formatterInstance, "ReTuRn TrUe")
				if got != tt.want {
					t.Fatalf("formatted output = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("options apply in order and ignore nil", func(t *testing.T) {
		formatterInstance := mustNewFormatter(
			t,
			WithPrintWidth(40),
			nil,
			WithPrintWidth(60),
		)

		if formatterInstance.config.PrintWidth != 60 {
			t.Fatalf("print width = %d, want 60", formatterInstance.config.PrintWidth)
		}
	})
}

func TestFormatterOptionValidation(t *testing.T) {
	tests := []struct {
		option Option
		name   string
		field  string
		value  string
	}{
		{
			name:   "zero print width",
			field:  "print width",
			value:  "0",
			option: WithPrintWidth(0),
		},
		{
			name:   "zero tab width",
			field:  "tab width",
			value:  "0",
			option: WithTabWidth(0),
		},
		{
			name:   "unsupported case mode",
			field:  "case mode",
			value:  "99",
			option: WithCaseMode(CaseMode(99)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatterInstance, err := New(tt.option)
			if err == nil {
				t.Fatal("New() error = nil, want validation error")
			}
			if formatterInstance != nil {
				t.Fatal("formatter != nil, want nil")
			}

			requireValidationError(t, err, tt.field, tt.value)
		})
	}
}

func TestFormatterOptionValidationAggregatesFailures(t *testing.T) {
	formatterInstance, err := New(
		WithPrintWidth(0),
		WithTabWidth(0),
	)
	if err == nil {
		t.Fatal("New() error = nil, want validation errors")
	}
	if formatterInstance != nil {
		t.Fatal("formatter != nil, want nil")
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("New() error = %T, want joined errors", err)
	}

	branches := joined.Unwrap()
	if len(branches) != 2 {
		t.Fatalf("joined errors = %d, want 2", len(branches))
	}

	requireValidationError(t, branches[0], "print width", "0")
	requireValidationError(t, branches[1], "tab width", "0")
}

func TestNewReturnsOptionError(t *testing.T) {
	reason := errors.New("test option failure")
	want := options.ValidationError{
		Field:  "formatter",
		Value:  "invalid",
		Reason: reason,
	}
	invalid := func(_ *internal.Config) error {
		return want
	}

	formatterInstance, err := New(invalid)
	if err == nil {
		t.Fatal("New() error = nil, want option error")
	}
	if formatterInstance != nil {
		t.Fatal("formatter != nil, want nil")
	}
	if !errors.Is(err, reason) {
		t.Fatalf("New() error = %v, want reason %v", err, reason)
	}

	requireValidationError(t, err, want.Field, want.Value)
}

func TestFormatterConcurrentUse(t *testing.T) {
	formatterInstance := mustNewFormatter(
		t,
		WithPrintWidth(40),
		WithTabWidth(2),
		WithSingleQuote(true),
	)
	input := `RETURN FOR value IN ["one", "two"] { RETURN { value: value } }`
	want := "return for value in ['one', 'two'] {\n  return { value: value }\n}"

	const workers = 32

	errs := make(chan error, workers)
	var group sync.WaitGroup
	group.Add(workers)

	for range workers {
		go func() {
			defer group.Done()

			var output bytes.Buffer
			if err := formatterInstance.Format(&output, source.NewAnonymous(input)); err != nil {
				errs <- err

				return
			}
			if output.String() != want {
				errs <- fmt.Errorf("formatted output = %q, want %q", output.String(), want)
			}
		}()
	}

	group.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

func formatForOptionTest(t testing.TB, formatterInstance *Formatter, input string) string {
	t.Helper()

	var output bytes.Buffer
	if err := formatterInstance.Format(&output, source.NewAnonymous(input)); err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	return output.String()
}

func requireValidationError(t testing.TB, err error, field, value string) {
	t.Helper()

	var got options.ValidationError
	if !errors.As(err, &got) {
		t.Fatalf("error = %T, want options.ValidationError", err)
	}
	if got.Field != field || got.Value != value {
		t.Fatalf("validation error = %+v, want field %q and value %q", got, field, value)
	}
}
