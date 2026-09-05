package stdlib_test

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestStringEncodingCryptoCapabilities(t *testing.T) {
	groups := map[stdlib.Group][]string{
		stdlib.Strings:  {"LOWER", "UPPER", "TRIM", "LTRIM", "RTRIM", "CONTAINS", "STARTS_WITH", "ENDS_WITH", "FIND_FIRST", "FIND_LAST", "SUBSTRING", "LEFT", "RIGHT", "SPLIT", "JOIN", "CONCAT", "REPLACE", "REPEAT", "FMT", "LIKE", "REGEX_TEST", "REGEX_FIND", "REGEX_FIND_ALL", "REGEX_REPLACE", "REGEX_SPLIT"},
		stdlib.Encoding: {"ENCODING::JSON_PARSE", "ENCODING::JSON_STRINGIFY", "ENCODING::QUERY_ESCAPE", "ENCODING::QUERY_UNESCAPE", "ENCODING::BASE64_ENCODE", "ENCODING::BASE64_DECODE", "ENCODING::HTML_ESCAPE", "ENCODING::HTML_UNESCAPE"},
		stdlib.Crypto:   {"CRYPTO::MD5", "CRYPTO::SHA1", "CRYPTO::SHA256", "CRYPTO::SHA512", "CRYPTO::RANDOM_TOKEN"},
	}
	for group, want := range groups {
		t.Run(string(group), func(t *testing.T) {
			for index, name := range want {
				want[index] = strings.ToLower(name)
			}

			got := functionNames(buildFunctions(t, stdlib.Only(group)))
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("functions = %v, want %v", got, want)
			}

			for _, set := range []stdlib.Set{stdlib.Full(), stdlib.Safe()} {
				functions := buildFunctions(t, set)
				without := buildFunctions(t, set.Without(group))
				for _, name := range want {
					if !functions.Has(name) || without.Has(name) {
						t.Fatalf("capability selection failed for %s", name)
					}

					if group != stdlib.Strings && (!functions.A1().Has(name) || functions.Var().Has(name)) {
						t.Fatalf("%s must have only a unary registration", name)
					}
				}
			}
		})
	}

	functions := buildFunctions(t, stdlib.Full())
	for _, name := range []string{"JSON_PARSE", "JSON_STRINGIFY", "ENCODE_URI_COMPONENT", "DECODE_URI_COMPONENT", "TO_BASE64", "FROM_BASE64", "ESCAPE_HTML", "UNESCAPE_HTML", "MD5", "SHA1", "SHA512", "RANDOM_TOKEN", "CONCAT_SEPARATOR", "SUBSTITUTE", "REGEX_MATCH"} {
		if functions.Has(name) {
			t.Fatalf("obsolete global %s remains registered", name)
		}
	}

	path := buildFunctions(t, stdlib.Only(stdlib.Path))
	if path.Has("JOIN") || !path.Var().Has("PATH::JOIN") {
		t.Fatal("path joining must be namespaced")
	}

	if !functions.A2().Has("JOIN") || functions.Var().Has("JOIN") || !functions.Var().Has("PATH::JOIN") {
		t.Fatal("string and path join registrations must coexist")
	}
}
