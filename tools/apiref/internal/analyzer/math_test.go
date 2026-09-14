package analyzer_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertMathMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()

	signatures := map[apicatalog.FunctionRef][]api.Signature{}
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			signatures[apicatalog.FunctionRef{Namespace: namespace.Name, Name: function.Name}] = function.Signatures
			if namespace.Name == "math" {
				for _, signature := range function.Signatures {
					if signature.Deprecated != "" {
						t.Fatalf("canonical math::%s is deprecated", function.Name)
					}
				}
			}
		}
	}

	for _, name := range strings.Fields("pi abs acos asin atan ceil cos degrees exp exp2 floor log log2 log10 max mean median min radians round sin sqrt stddev stddev_sample sum tan variance variance_sample atan2 pow percentile") {
		arity := 1
		if name == "pi" {
			arity = 0
		}

		if name == "atan2" || name == "pow" || name == "percentile" {
			arity = 2
		}

		canonical := signatures[apicatalog.FunctionRef{Namespace: "math", Name: name}]
		if len(canonical) != 1 || !hasSignature(canonical, arity, false) {
			t.Fatalf("math::%s signature = %+v", name, canonical)
		}

		global := name
		switch name {
		case "mean":
			global = "average"
		case "variance":
			global = "variance_population"
		case "stddev":
			global = "stddev_population"
		}

		legacy := signatures[apicatalog.FunctionRef{Name: global}]
		count := 1
		if name == "percentile" {
			count = 2
		}

		if len(legacy) != count {
			t.Fatalf("%s legacy signatures = %+v", global, legacy)
		}

		for _, signature := range legacy {
			if signature.Deprecated != "Use math::"+name+" instead." {
				t.Fatalf("%s deprecation = %q", global, signature.Deprecated)
			}
		}
	}

	for _, name := range []string{"average", "variance_population", "stddev_population", "range", "rand"} {
		if _, ok := signatures[apicatalog.FunctionRef{Namespace: "math", Name: name}]; ok {
			t.Fatalf("unexpected math::%s", name)
		}
	}

	for _, signature := range signatures[apicatalog.FunctionRef{Name: "rand"}] {
		if !strings.Contains(signature.Deprecated, "planned") || !strings.Contains(signature.Deprecated, "no replacement") {
			t.Fatalf("rand deprecation = %q", signature.Deprecated)
		}
	}

	rangeSignatures := signatures[apicatalog.FunctionRef{Name: "range"}]
	if len(rangeSignatures) != 2 {
		t.Fatal("range requires two overloads")
	}

	for _, signature := range rangeSignatures {
		if signature.Deprecated != "Use arrays::range instead." {
			t.Fatalf("range deprecation = %q", signature.Deprecated)
		}
	}

	categories := map[apicatalog.FunctionRef]string{}
	for _, category := range catalog.Categories {
		for _, function := range category.Functions {
			categories[function] = category.ID
		}
	}

	for _, namespace := range []string{"", "arrays"} {
		if categories[apicatalog.FunctionRef{Namespace: namespace, Name: "range"}] != "arrays" {
			t.Fatal("both range forms must be in Arrays")
		}
	}
}
