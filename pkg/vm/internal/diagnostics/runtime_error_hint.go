package diagnostics

import (
	"fmt"
	"regexp"
	"strings"

	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

type arityExpectation struct {
	CallTarget string
	Expected   string
	AtLeast    bool
}

var (
	runtimeArityCallablePattern  = regexp.MustCompile(`^(.+?) expects (.+) argument[s]?, (?:but )?got \d+$`)
	runtimeArityExpectedPattern  = regexp.MustCompile(`^expected number of arguments (.+), but got \d+$`)
	runtimeArityAtLeastPattern   = regexp.MustCompile(`^expected at least (\d+) argument[s]?, but got \d+$`)
	runtimeArgumentTypePattern   = regexp.MustCompile(`^argument (\d+) expects (.+), but got .+$`)
	runtimeExpectedTypePattern   = regexp.MustCompile(`^expected (.+), but got .+$`)
	runtimeArgumentMustBePattern = regexp.MustCompile(`^argument (\d+) must be (.+)$`)
	runtimeMustBePattern         = regexp.MustCompile(`^must be (.+)$`)
)

func synthesizeRuntimeHint(spec runtimeDiagnosticSpec) string {
	switch spec.Kind {
	case ArityError:
		return synthesizeArityHint(spec.Note)
	case InvalidArgument:
		return synthesizeInvalidArgumentHint(spec.Note)
	case pkgdiagnostics.TypeError:
		return synthesizeTypeHint(spec.Message, spec.Note)
	default:
		return ""
	}
}

func arityLabel(note string) string {
	expectation, ok := parseArityExpectation(note)
	if !ok || expectation.CallTarget == "" {
		return "wrong number of arguments"
	}

	return fmt.Sprintf("wrong number of arguments in call to %s", expectation.CallTarget)
}

func synthesizeArityHint(note string) string {
	expectation, ok := parseArityExpectation(note)
	if !ok {
		return ""
	}

	target := "this call"
	if expectation.CallTarget != "" {
		target = expectation.CallTarget
	}

	if expectation.AtLeast {
		return fmt.Sprintf(
			"Pass at least %s %s to %s",
			expectation.Expected,
			argumentWord(expectation.Expected),
			target,
		)
	}

	return fmt.Sprintf(
		"Pass %s %s to %s",
		expectation.Expected,
		argumentWord(expectation.Expected),
		target,
	)
}

func synthesizeTypeHint(message string, note string) string {
	if match := runtimeArgumentTypePattern.FindStringSubmatch(strings.TrimSpace(note)); len(match) == 3 {
		return fmt.Sprintf("Convert argument %s to %s before this call", match[1], match[2])
	}

	if match := runtimeExpectedTypePattern.FindStringSubmatch(strings.TrimSpace(note)); len(match) == 2 {
		if message == "invalid argument type" {
			return fmt.Sprintf("Pass a value of type %s to this call", match[1])
		}

		return fmt.Sprintf("Convert the value to %s before this operation", match[1])
	}

	return ""
}

func synthesizeInvalidArgumentHint(note string) string {
	trimmed := strings.TrimSpace(note)

	if match := runtimeArgumentMustBePattern.FindStringSubmatch(trimmed); len(match) == 3 {
		return fmt.Sprintf("Pass argument %s with a value that is %s", match[1], match[2])
	}

	if match := runtimeMustBePattern.FindStringSubmatch(trimmed); len(match) == 2 {
		return fmt.Sprintf("Pass a value that is %s", match[1])
	}

	return ""
}

func parseArityExpectation(note string) (arityExpectation, bool) {
	trimmed := strings.TrimSpace(note)

	if match := runtimeArityAtLeastPattern.FindStringSubmatch(trimmed); len(match) == 2 {
		return arityExpectation{
			Expected: match[1],
			AtLeast:  true,
		}, true
	}

	if match := runtimeArityExpectedPattern.FindStringSubmatch(trimmed); len(match) == 2 {
		return arityExpectation{
			Expected: strings.TrimSpace(match[1]),
		}, true
	}

	if match := runtimeArityCallablePattern.FindStringSubmatch(trimmed); len(match) == 3 {
		return arityExpectation{
			CallTarget: normalizeCallTarget(match[1]),
			Expected:   strings.TrimSpace(match[2]),
		}, true
	}

	return arityExpectation{}, false
}

func normalizeCallTarget(subject string) string {
	subject = strings.TrimSpace(subject)

	switch {
	case strings.HasPrefix(subject, "UDF '") && strings.HasSuffix(subject, "'"):
		return strings.TrimSuffix(strings.TrimPrefix(subject, "UDF '"), "'")
	case strings.HasPrefix(subject, "function '") && strings.HasSuffix(subject, "'"):
		return strings.TrimSuffix(strings.TrimPrefix(subject, "function '"), "'")
	case strings.HasPrefix(subject, "'") && strings.HasSuffix(subject, "'"):
		return strings.TrimSuffix(strings.TrimPrefix(subject, "'"), "'")
	default:
		return subject
	}
}

func argumentWord(expected string) string {
	if strings.TrimSpace(expected) == "1" {
		return "argument"
	}

	return "arguments"
}
