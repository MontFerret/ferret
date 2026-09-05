package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestModernStringContracts(t *testing.T) {
	specs := []spec.Spec{
		S(`RETURN LEFT("😀abc",1)`, "😀", "left uses rune bounds"),
		S(`RETURN RIGHT("abc😀",5)`, "abc😀", "right clamps rune bounds"),
		S(`RETURN SUBSTRING("😀é文",1,9223372036854775807)`, "é文", "substring length cannot overflow"),
		Array(`RETURN [TRIM("\t\n x \n\t"), LTRIM("\t\n x"), RTRIM("x \n\t")]`, []any{"x", "x", "x"}, "Unicode whitespace defaults"),
		Array(`RETURN [CONTAINS("😀é","é"),STARTS_WITH("😀é","😀"),ENDS_WITH("😀é","é")]`, []any{true, true, true}, "Boolean string predicates"),
		S(`RETURN JOIN(["","é","😀",""],"|")`, "|é|😀|", "join preserves empty elements"),
		S(`RETURN PATH::JOIN("a","b")`, "a/b", "path join remains separate"),
		S(`RETURN REPEAT("😀",3)`, "😀😀😀", "repeat Unicode"),
		S(`RETURN REPEAT("a",9223372036854775807) ON ERROR RETURN "caught"`, "caught", "oversized ASCII repetition is catchable"),
		S(`RETURN REPEAT("😀",16777217) ON ERROR RETURN "caught"`, "caught", "repetition byte limit is catchable"),
		Array(`RETURN [REPLACE("aaa","a","x",2),REPLACE("aaa","a","",0)]`, []any{"xxa", "aaa"}, "literal replacement limits"),
		Array(`RETURN SPLIT("a,b,c",",",2)`, []any{"a", "b,c"}, "split remainder"),
		Array(`RETURN REGEX_SPLIT("a1b2c","[0-9]",2)`, []any{"a", "b2c"}, "regex split remainder"),
		S(`RETURN FMT("{{{1}}} {0} {1}","é",42)`, "{42} é 42", "indexed formatting and escaped braces"),
		S(`RETURN LIKE("a_b%","a_b%") AND NOT LIKE("acb","a_b")`, true, "glob treats SQL characters literally"),
		S(`RETURN REGEX_TEST("ABC","(?i)abc")`, true, "regex inline flags"),
		S(`RETURN REGEX_FIND("é😀","(?P<letter>é)(😀)?") == {match:"é😀",groups:["é","😀"],named:{letter:"é"}}`, true, "structured regex captures"),
		Nil(`RETURN REGEX_FIND("a","z")`, "no regex match returns None"),
		S(`RETURN LENGTH(REGEX_FIND_ALL("aba","a"))`, 2, "regex finds all matches"),
		Array(`RETURN REGEX_FIND_ALL("a","z")`, []any{}, "no matches returns empty array"),
		S(`RETURN REGEX_REPLACE("a1 a2","(?P<letter>a)([0-9])","${letter}:$2:$$")`, "a:1:$ a:2:$", "regex replacement expansion"),
		S(`RETURN ENCODING::JSON_PARSE(ENCODING::JSON_STRINGIFY({x:[1,NONE,true]})) == {x:[1,NONE,true]}`, true, "JSON namespace preserves values"),
		S(`RETURN ENCODING::QUERY_ESCAPE("a +/?&é")`, "a+%2B%2F%3F%26%C3%A9", "query-form encoding"),
		S(`RETURN ENCODING::QUERY_UNESCAPE("a+%2B")`, "a +", "query-form decoding"),
		S(`RETURN IS_BINARY(ENCODING::BASE64_DECODE("AP+A"))`, true, "Base64 decoding returns Binary"),
		S(`RETURN ENCODING::BASE64_ENCODE(ENCODING::BASE64_DECODE("AP+A"))`, "AP+A", "Base64 preserves arbitrary bytes"),
		S(`RETURN ENCODING::HTML_UNESCAPE(ENCODING::HTML_ESCAPE("<é>"))`, "<é>", "HTML escaping namespace"),
		S(`RETURN CRYPTO::SHA256("foobar")`, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", "SHA256 vector"),
		S(`LET token = CRYPTO::RANDOM_TOKEN(32) RETURN LENGTH(token) == 32 AND REGEX_TEST(token,"^[a-zA-Z0-9]+$")`, true, "secure token contract"),
	}
	for _, test := range []struct {
		query   string
		message string
	}{
		{query: `RETURN LEFT("😀",-1)`, message: "invalid argument"},
		{query: `RETURN RIGHT("😀",-1)`, message: "invalid argument"},
		{query: `RETURN REPEAT("😀",-1)`, message: "invalid argument"},
		{query: `RETURN REPEAT("😀",9223372036854775807)`, message: "invalid argument"},
		{query: `RETURN REPEAT("a",9223372036854775807)`, message: "invalid argument"},
		{query: `RETURN REPEAT("a",67108865)`, message: "invalid argument"},
		{query: `RETURN REPEAT("😀",16777217)`, message: "invalid argument"},
		{query: `RETURN SUBSTRING("x",0,true)`, message: "invalid argument type"},
		{query: `RETURN SPLIT("a",",",-1)`, message: "invalid argument"},
		{query: `RETURN REPLACE("a","a","x",-1)`, message: "invalid argument"},
		{query: `RETURN REGEX_SPLIT("a",",",true)`, message: "invalid argument type"},
		{query: `RETURN CONTAINS("123",1)`, message: "invalid argument type"},
		{query: `RETURN LOWER(1)`, message: "invalid argument type"},
		{query: `RETURN TRIM("x",NONE)`, message: "invalid argument type"},
		{query: `RETURN FIND_FIRST("é","é",true)`, message: "invalid argument type"},
		{query: `RETURN FIND_LAST("é","é",0,false)`, message: "invalid argument type"},
		{query: `RETURN JOIN(["x",NONE],",")`, message: "invalid argument"},
		{query: `RETURN LIKE("x","x",1)`, message: "invalid argument type"},
		{query: `RETURN FMT("{0} {}",1)`, message: "invalid argument"},
		{query: `RETURN FMT("literal",1)`, message: "invalid argument"},
		{query: `RETURN FMT("{1}",1,2)`, message: "invalid argument"},
		{query: `RETURN REGEX_FIND("x","[")`, message: "invalid argument"},
		{query: `RETURN ENCODING::JSON_PARSE(1)`, message: "invalid argument type"},
		{query: `RETURN ENCODING::BASE64_DECODE("!")`, message: "invalid argument"},
		{query: `RETURN CRYPTO::SHA256(1)`, message: "invalid argument type"},
		{query: `RETURN CRYPTO::RANDOM_TOKEN(0)`, message: "invalid argument"},
	} {
		specs = append(specs, spec.NewSpec(test.query, test.query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: test.message}))
	}

	for _, query := range []string{
		`RETURN CONTAINS("a","a",true)`, `RETURN REGEX_TEST("a","a",true)`,
		`RETURN REGEX_REPLACE("a","a","b",true)`, `RETURN REGEX_SPLIT("a",",",2,true)`,
		`RETURN REPLACE("a","a")`, `RETURN JOIN(["a"],",","extra")`,
	} {
		specs = append(specs, spec.NewSpec(query, query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid number of arguments"}))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, level.String(), c, specs)
	}
}
