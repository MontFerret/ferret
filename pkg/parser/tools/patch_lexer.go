package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := filepath.Join("fql", "fql_lexer.go")
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("read %s: %v", path, err)
	}

	if bytes.Contains(data, []byte("templateDepth []int")) {
		return
	}

	lines := strings.Split(string(data), "\n")
	insertAt := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "type FqlLexer struct {") {
			for j := i + 1; j < len(lines); j++ {
				if strings.Contains(lines[j], "modeNames") {
					insertAt = j + 1
					break
				}
				if strings.Contains(lines[j], "}") {
					insertAt = j
					break
				}
			}
			if insertAt == -1 {
				insertAt = i + 1
			}
			break
		}
	}

	if insertAt == -1 {
		fatalf("could not locate FqlLexer struct in %s", path)
	}

	lines = append(lines[:insertAt], append([]string{"\ttemplateDepth []int"}, lines[insertAt:]...)...)

	out := strings.Join(lines, "\n")
	formatted, err := format.Source([]byte(out))
	if err != nil {
		if writeErr := os.WriteFile(path, []byte(out), 0o644); writeErr != nil {
			fatalf("write %s after format error: %v", path, writeErr)
		}
		fatalf("format %s: %v", path, err)
	}

	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		fatalf("write %s: %v", path, err)
	}
}

func fatalf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
