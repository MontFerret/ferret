package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type target struct {
	path string
	skip map[int]bool
}

func main() {
	if len(os.Args) != 3 {
		panic("usage: extract WEBSITE OUTPUT")
	}

	website := os.Args[1]
	output := os.Args[2]
	targets := []target{
		{path: "content/docs/embedding/go/application-guide.md"},
		{path: "content/docs/embedding/go/custom-functions.md"},
		{path: "content/docs/embedding/go/getting-started.md"},
		{path: "content/docs/embedding/go/host-values.md"},
		{path: "content/docs/embedding/go/migrating-from-v1.md", skip: map[int]bool{1: true}},
		{path: "content/docs/embedding/go/parameters.md"},
		{path: "content/docs/embedding/go/programs.md"},
		{path: "content/docs/embedding/go/value-encoders.md"},
		{path: "content/docs/guides/existing-html.md"},
		{path: "content/docs/guides/precompiled-programs.md"},
	}

	for _, item := range targets {
		if err := extract(website, output, item); err != nil {
			panic(err)
		}
	}
}

func extract(website, output string, item target) error {
	file, err := os.Open(filepath.Join(website, item.path))
	if err != nil {
		return err
	}
	defer file.Close()

	var block strings.Builder
	inGoBlock := false
	program := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == `{{< code lang="go" >}}`:
			inGoBlock = true
			block.Reset()
		case line == `{{</ code >}}` && inGoBlock:
			inGoBlock = false
			code := block.String()
			if !strings.HasPrefix(code, "package main\n") {
				continue
			}

			program++
			if item.skip[program] {
				continue
			}

			name := strings.TrimSuffix(filepath.Base(item.path), filepath.Ext(item.path))
			dir := filepath.Join(output, fmt.Sprintf("%s-%d", name, program))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
				return err
			}

			for _, codeLine := range strings.Split(code, "\n") {
				pattern, found := strings.CutPrefix(strings.TrimSpace(codeLine), "//go:embed ")
				if !found || strings.ContainsAny(pattern, "*?[") {
					continue
				}

				asset := filepath.Join(dir, pattern)
				if err := os.MkdirAll(filepath.Dir(asset), 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(asset, nil, 0o644); err != nil {
					return err
				}
			}
		case inGoBlock:
			block.WriteString(line)
			block.WriteByte('\n')
		}
	}

	return scanner.Err()
}
