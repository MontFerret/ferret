package core

import (
	"regexp"
	"strings"
)

func ExplainSyntaxError(err string) (msg string, hint string) {
	var matched bool
	parsers := []func(string) (string, string, bool){
		explainNoViableAltError,
		explainExtraneousError,
	}

	for _, parser := range parsers {
		msg, hint, matched = parser(err)

		if matched {
			return
		}
	}

	msg = "Syntax error"
	hint = "Check the syntax of your code. It may be missing a keyword, operator, or punctuation."

	return
}

func explainExtraneousError(err string) (msg string, hint string, matched bool) {
	recognizer := regexp.MustCompile("extraneous input '<EOF>' expecting")

	if !recognizer.MatchString(err) {
		return "", "", false
	}

	return "Extraneous input at end of file", "Check the syntax of your code. It may be missing a keyword, operator, or punctuation", true
}

func explainNoViableAltError(err string) (msg string, hint string, matched bool) {
	recognizer := regexp.MustCompile("no viable alternative at input '(\\w+).+'")

	matches := recognizer.FindAllStringSubmatch(err, -1)

	if len(matches) == 0 {
		return "", "", false
	}

	keyword := matches[0][1]

	switch strings.ToLower(keyword) {
	case "return":
		msg = "Unexpected 'return' keyword"
		hint = "Did you mean to return a value?"
	}

	return
}
