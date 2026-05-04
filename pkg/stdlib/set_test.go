package stdlib_test

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func buildFunctions(t *testing.T, set stdlib.Set) *runtime.Functions {
	t.Helper()

	ns := runtime.NewLibrary()
	if err := set.Register(ns); err != nil {
		t.Fatalf("failed to register stdlib set: %v", err)
	}

	funcs, err := ns.Build()
	if err != nil {
		t.Fatalf("failed to build functions: %v", err)
	}

	return funcs
}

func functionNames(funcs *runtime.Functions) []string {
	names := funcs.List()
	sort.Strings(names)

	return names
}

func TestFullRegistersRepresentativeFunctions(t *testing.T) {
	t.Parallel()

	funcs := buildFunctions(t, stdlib.Full())

	for _, name := range []string{
		"CONCAT",
		"IO::FS::READ",
		"IO::NET::HTTP::GET",
	} {
		if !funcs.Has(name) {
			t.Fatalf("expected full stdlib to register %s", name)
		}
	}
}

func TestSafeExcludesExternalIO(t *testing.T) {
	t.Parallel()

	funcs := buildFunctions(t, stdlib.Safe())

	if !funcs.Has("CONCAT") {
		t.Fatal("expected safe stdlib to keep non-IO functions")
	}

	for _, name := range []string{
		"IO::FS::READ",
		"IO::FS::WRITE",
		"IO::NET::HTTP::GET",
		"IO::NET::HTTP::POST",
	} {
		if funcs.Has(name) {
			t.Fatalf("expected safe stdlib to exclude %s", name)
		}
	}
}

func TestSafeMatchesFullWithoutIO(t *testing.T) {
	t.Parallel()

	safeNames := functionNames(buildFunctions(t, stdlib.Safe()))
	withoutIONames := functionNames(buildFunctions(t, stdlib.Full().Without(stdlib.IO)))

	if !reflect.DeepEqual(safeNames, withoutIONames) {
		t.Fatalf("expected Safe to match Full().Without(IO), got %v vs %v", safeNames, withoutIONames)
	}
}

func TestEmptyRegistersNoFunctions(t *testing.T) {
	t.Parallel()

	funcs := buildFunctions(t, stdlib.Empty())

	if funcs.Size() != 0 {
		t.Fatalf("expected empty stdlib to register no functions, got %d", funcs.Size())
	}
}

func TestOnlyRegistersNestedIOGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set      stdlib.Set
		included []string
		excluded []string
	}{
		{
			name:     "fs",
			set:      stdlib.Only(stdlib.FS),
			included: []string{"IO::FS::READ", "IO::FS::WRITE"},
			excluded: []string{"IO::NET::HTTP::GET", "CONCAT"},
		},
		{
			name:     "net",
			set:      stdlib.Only(stdlib.NET),
			included: []string{"IO::NET::HTTP::GET", "IO::NET::HTTP::POST"},
			excluded: []string{"IO::FS::READ", "CONCAT"},
		},
		{
			name:     "io",
			set:      stdlib.Only(stdlib.IO),
			included: []string{"IO::FS::READ", "IO::NET::HTTP::GET"},
			excluded: []string{"CONCAT"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			funcs := buildFunctions(t, tt.set)

			for _, name := range tt.included {
				if !funcs.Has(name) {
					t.Fatalf("expected %s to be registered", name)
				}
			}

			for _, name := range tt.excluded {
				if funcs.Has(name) {
					t.Fatalf("expected %s not to be registered", name)
				}
			}
		})
	}
}

func TestSetModifiersAreImmutable(t *testing.T) {
	t.Parallel()

	base := stdlib.Only(stdlib.Strings)
	withFS := base.With(stdlib.FS)
	withoutStrings := base.Without(stdlib.Strings)

	baseFuncs := buildFunctions(t, base)
	if !baseFuncs.Has("CONCAT") {
		t.Fatal("expected base set to retain strings")
	}
	if baseFuncs.Has("IO::FS::READ") {
		t.Fatal("expected base set not to gain FS from derived set")
	}

	withFSFuncs := buildFunctions(t, withFS)
	if !withFSFuncs.Has("CONCAT") || !withFSFuncs.Has("IO::FS::READ") {
		t.Fatal("expected With-derived set to contain strings and FS")
	}

	withoutStringsFuncs := buildFunctions(t, withoutStrings)
	if withoutStringsFuncs.Has("CONCAT") {
		t.Fatal("expected Without-derived set to remove strings")
	}
}

func TestSetRegisterRejectsInvalidGroup(t *testing.T) {
	t.Parallel()

	ns := runtime.NewLibrary()
	err := stdlib.Only(stdlib.Group("unknown")).Register(ns)
	if err == nil {
		t.Fatal("expected invalid group to fail registration")
	}

	if !strings.Contains(err.Error(), "invalid stdlib group(s): unknown") {
		t.Fatalf("expected invalid group error, got: %v", err)
	}
}

func TestZeroValueSetRegistersNoFunctions(t *testing.T) {
	t.Parallel()

	funcs := buildFunctions(t, stdlib.Set{})

	if funcs.Size() != 0 {
		t.Fatalf("expected zero-value set to register no functions, got %d", funcs.Size())
	}
}
