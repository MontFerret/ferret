package runtime_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRegexpRuntimeContract(t *testing.T) {
	value, err := runtime.NewRegexp(`^<item>&[a-z]+>$`)
	if err != nil {
		t.Fatalf("compile regexp: %v", err)
	}

	if value.Type() != runtime.TypeRegexp || runtime.TypeOf(value) != runtime.TypeRegexp {
		t.Fatalf("regexp type = %s, want %s", runtime.TypeOf(value), runtime.TypeRegexp)
	}

	if got := runtime.TypeName(value.Type()); got != "Regexp" {
		t.Fatalf("regexp type name = %q, want Regexp", got)
	}

	if got, want := value.String(), `^<item>&[a-z]+>$`; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if !value.Match(runtime.NewString("<item>&value>")) || !value.MatchString("<item>&value>") {
		t.Fatal("regexp did not match the expected value")
	}

	if value.MatchString("item&value") {
		t.Fatal("regexp matched an invalid value")
	}

	copy, ok := value.Copy().(*runtime.Regexp)
	if !ok {
		t.Fatalf("Copy() type = %T, want *runtime.Regexp", value.Copy())
	}

	if copy == value {
		t.Fatal("Copy() returned the original regexp")
	}

	equal, err := value.Equal(t.Context(), copy)
	if err != nil || !equal {
		t.Fatalf("Equal(copy) = %v, %v, want true, nil", equal, err)
	}

	if value.Hash() != copy.Hash() {
		t.Fatal("equal regexps have different hashes")
	}

	later, err := runtime.NewRegexp(`^z`)
	if err != nil {
		t.Fatalf("compile later regexp: %v", err)
	}

	ordering, err := value.Compare(t.Context(), later)
	if err != nil || ordering != runtime.Less {
		t.Fatalf("Compare() = %v, %v, want Less, nil", ordering, err)
	}

	equal, err = value.Equal(t.Context(), runtime.NewString(value.String()))
	if err != nil || equal {
		t.Fatalf("Equal(String) = %v, %v, want false, nil", equal, err)
	}

	if _, err := value.Compare(t.Context(), runtime.NewString(value.String())); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Compare(String) error = %v, want ErrInvalidOperation", err)
	}
}

func TestNewRegexpRejectsInvalidPattern(t *testing.T) {
	if _, err := runtime.NewRegexp("["); err == nil {
		t.Fatal("NewRegexp accepted an invalid pattern")
	}
}
