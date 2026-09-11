package runtime_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCopyConcreteArray(t *testing.T) {
	nested := runtime.NewArrayWith(runtime.Int(1))
	source := runtime.NewArrayWith(nested, runtime.Int(2))
	copied, err := runtime.Copy(source)
	if err != nil {
		t.Fatal(err)
	}

	var _ *runtime.Array = copied
	if copied == source || copied.String() != source.String() {
		t.Fatalf("copy = %v, source = %v", copied, source)
	}

	first, err := copied.At(t.Context(), 0)
	if err != nil || first != nested {
		t.Fatalf("shallow element = %v, error = %v", first, err)
	}

	if err := copied.SetAt(t.Context(), 1, runtime.Int(3)); err != nil {
		t.Fatal(err)
	}

	if source.String() != "[[1],2]" || copied.String() != "[[1],3]" {
		t.Fatalf("replacement shared storage: source = %v, copy = %v", source, copied)
	}
}

func TestCopyList(t *testing.T) {
	var source runtime.List = runtime.NewArrayWith(runtime.Int(1))
	copied, err := runtime.Copy(source)
	if err != nil {
		t.Fatal(err)
	}

	var _ runtime.List = copied
	if copied == source || copied.String() != "[1]" {
		t.Fatalf("copy = %v, source = %v", copied, source)
	}

	host := &copyContractList{List: source, copied: copied}
	fromHost, err := runtime.Copy[runtime.List](host)
	if err != nil || fromHost != copied || host.calls != 1 {
		t.Fatalf("host copy = %v, error = %v, calls = %d", fromHost, err, host.calls)
	}
}

func TestCopyTypeMismatch(t *testing.T) {
	for _, returned := range []runtime.Value{runtime.True, nil} {
		t.Run(fmt.Sprintf("%T", returned), func(t *testing.T) {
			host := &copyContractList{copied: returned}
			copied, err := runtime.Copy[runtime.List](host)
			if copied != nil || !errors.Is(err, runtime.ErrInvalidType) || host.calls != 1 {
				t.Fatalf("copy = %v, error = %v, calls = %d", copied, err, host.calls)
			}

			for _, detail := range []string{"expected List", "copy of *runtime_test.copyContractList", fmt.Sprintf("returned %T", returned)} {
				if !strings.Contains(err.Error(), detail) {
					t.Errorf("error %q missing %q", err, detail)
				}
			}
		})
	}

	t.Run("concrete", func(t *testing.T) {
		host := &copyContractList{copied: runtime.EmptyArray()}
		copied, err := runtime.Copy(host)
		if copied != nil || !errors.Is(err, runtime.ErrInvalidType) || host.calls != 1 {
			t.Fatalf("copy = %v, error = %v, calls = %d", copied, err, host.calls)
		}

		if !strings.Contains(err.Error(), "returned *runtime.Array") {
			t.Fatalf("missing returned concrete type: %v", err)
		}
	})
}

func TestCopyPreservesOnlyStaticType(t *testing.T) {
	host := &copyContractList{copied: runtime.True}
	var source runtime.Value = host
	copied, err := runtime.Copy(source)
	if err != nil || copied != runtime.True || host.calls != 1 {
		t.Fatalf("Value copy = %v, error = %v, calls = %d", copied, err, host.calls)
	}
}

func TestToListCopyContract(t *testing.T) {
	host := &copyContractList{List: runtime.NewArrayWith(runtime.Int(1)), copied: runtime.True}
	copied, err := runtime.ToList(t.Context(), host)
	if copied != nil || !errors.Is(err, runtime.ErrInvalidType) || host.calls != 1 {
		t.Fatalf("copy = %v, error = %v, calls = %d", copied, err, host.calls)
	}

	if host.String() != "[1]" {
		t.Fatalf("source changed: %v", host)
	}
}
