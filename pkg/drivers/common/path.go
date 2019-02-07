package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"strings"
)

func PathToString(path []core.Value) string {
	spath := make([]string, 0, len(path))

	for i, s := range path {
		spath[i] = s.String()
	}

	return strings.Join(spath, ".")
}
