package format

import (
	"fmt"
	"strings"
	"testing"
)

func PrintDebug(t *testing.T, name string, input string, output strings.Builder) {
	t.Helper()
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Suite: %s\n", name))
	b.WriteString(fmt.Sprintf("Input Expression: %s\n", input))
	b.WriteString(fmt.Sprintf("Output Expression: %s\n", output))

	t.Log(b.String())
}
