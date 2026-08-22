package testing

import (
	"fmt"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func appendObjectPath(path, key string) string {
	if isDottedPathIdentifier(key) {
		return path + "." + key
	}

	return path + "[" + strconv.Quote(key) + "]"
}

func appendIndexPath(path string, index runtime.Int) string {
	return fmt.Sprintf("%s[%d]", path, index)
}

func isDottedPathIdentifier(value string) bool {
	if value == "" || !isASCIIPathLetter(value[0]) {
		return false
	}

	for index := 1; index < len(value); index++ {
		char := value[index]
		if !isASCIIPathLetter(char) && (char < '0' || char > '9') && char != '_' {
			return false
		}
	}

	return true
}

func isASCIIPathLetter(value byte) bool {
	return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z'
}
