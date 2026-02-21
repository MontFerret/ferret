package sdk

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	proxyErrorTemplate = "proxy target %T does not implement %s interface"
)

func ProxyError(target any, expected ...runtime.Type) error {
	if len(expected) == 0 {
		return runtime.Error(runtime.ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, "unknown"))
	}

	if len(expected) == 1 {
		return runtime.Error(runtime.ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, expected[0]))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = string(t)
	}

	expectedStr := strings.Join(strs, " or ")

	return runtime.Error(runtime.ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, expectedStr))
}
