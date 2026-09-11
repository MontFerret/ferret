package analyzer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertDateTimeMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()
	signatures := make(map[apicatalog.FunctionRef][]api.Signature)
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			signatures[apicatalog.FunctionRef{Namespace: namespace.Name, Name: function.Name}] = function.Signatures
			if namespace.Name == "datetime" {
				for _, signature := range function.Signatures {
					if signature.Deprecated != "" {
						t.Fatalf("canonical datetime::%s is deprecated", function.Name)
					}
				}
			}
		}
	}

	for legacy, canonical := range map[string]string{
		"date":               "parse",
		"date_add":           "add",
		"date_day":           "day",
		"date_dayofweek":     "day_of_week",
		"date_dayofyear":     "day_of_year",
		"date_days_in_month": "days_in_month",
		"date_format":        "format",
		"date_hour":          "hour",
		"date_leapyear":      "is_leap_year",
		"date_millisecond":   "millisecond",
		"date_minute":        "minute",
		"date_month":         "month",
		"date_quarter":       "quarter",
		"date_second":        "second",
		"date_subtract":      "subtract",
		"date_year":          "year",
		"now":                "now",
	} {
		old := signatures[apicatalog.FunctionRef{Name: legacy}]
		current := signatures[apicatalog.FunctionRef{Namespace: "datetime", Name: canonical}]
		if len(old) == 0 || len(old) != len(current) {
			t.Fatalf("%s / %s signature counts differ", legacy, canonical)
		}

		for index := range old {
			signature := old[index]
			if signature.Deprecated != "Use datetime::"+canonical+" instead." {
				t.Fatalf("%s deprecation = %q", legacy, signature.Deprecated)
			}

			signature.Deprecated = ""
			signature.Description = current[index].Description
			if !reflect.DeepEqual(signature, current[index]) {
				t.Fatalf("%s metadata differs from datetime::%s", legacy, canonical)
			}
		}
	}

	for legacy, canonical := range map[string]string{"date_compare": "same", "date_diff": "diff"} {
		old := signatures[apicatalog.FunctionRef{Name: legacy}]
		current := signatures[apicatalog.FunctionRef{Namespace: "datetime", Name: canonical}]
		if len(old) != 2 || len(current) != 1 || !hasSignature(old, 3, false) || !hasSignature(old, 4, false) || !hasSignature(current, 3, false) {
			t.Fatalf("%s / %s signatures are incorrect", legacy, canonical)
		}

		for _, signature := range old {
			if !strings.Contains(signature.Deprecated, "datetime::"+canonical) {
				t.Fatalf("%s lacks migration metadata", legacy)
			}
		}
	}

	diff := signatures[apicatalog.FunctionRef{Namespace: "datetime", Name: "diff"}][0]
	if diff.Return.Type.Kind != api.TypeKindNamed || diff.Return.Type.Name != "Float" {
		t.Fatalf("datetime::diff return type = %v", diff.Return.Type)
	}

	count := 0
	for _, category := range catalog.Categories {
		if category.ID != "datetime" {
			continue
		}

		for _, function := range category.Functions {
			if _, ok := signatures[function]; !ok {
				t.Fatalf("unknown datetime category entry %v", function)
			}

			count++
		}
	}

	if count != 38 {
		t.Fatalf("datetime catalog has %d functions, want 38", count)
	}
}
