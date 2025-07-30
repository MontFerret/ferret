package core

import (
	"regexp"
	"strings"
)

func NewSyntaxError(msg string, line, column int, offendingSymbol any) *CompilationError {
	err := &CompilationError{
		Message: msg,
		Hint:    "Check the syntax of your code. It may be missing a keyword, operator, or punctuation.",
	}

	var matched bool
	parsers := []func(*CompilationError, int, int, any) bool{
		parseNoViableAltError,
		parseExtraneousError,
	}

	for _, parser := range parsers {
		matched = parser(err, line, column, offendingSymbol)

		if matched {
			break
		}
	}

	return err
}

func parseExtraneousError(err *CompilationError, line, column int, offendingSymbol any) (matched bool) {
	recognizer := regexp.MustCompile("extraneous input '<EOF>' expecting")

	if !recognizer.MatchString(err.Message) {
		return false
	}

	err.Message = "Extraneous input at end of file"

	return true
}

func parseNoViableAltError(err *CompilationError, line, column int, offendingSymbol any) (matched bool) {
	recognizer := regexp.MustCompile("no viable alternative at input '(\\w+).+'")

	matches := recognizer.FindAllStringSubmatch(err.Message, -1)

	if len(matches) == 0 {
		return false
	}

	var msg, hint string
	keyword := matches[0][1]

	switch strings.ToLower(keyword) {
	case "return":
		msg = "Unexpected 'return' keyword"
		hint = "Did you mean to return a value?"
	}

	err.Message = msg
	err.Hint = hint

	return true
}
