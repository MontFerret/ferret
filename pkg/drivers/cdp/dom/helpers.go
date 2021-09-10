package dom

import (
	"bytes"
	"regexp"
	"strings"
)

var camelMatcher = regexp.MustCompile("[A-Za-z0-9]+")

func toCamelCase(input string) string {
	var buf bytes.Buffer

	matched := camelMatcher.FindAllString(input, -1)

	if matched == nil {
		return ""
	}

	for i, match := range matched {
		res := match

		if i > 0 {
			if len(match) > 1 {
				res = strings.ToUpper(match[0:1]) + match[1:]
			} else {
				res = strings.ToUpper(match)
			}
		}

		buf.WriteString(res)
	}

	return buf.String()
}
