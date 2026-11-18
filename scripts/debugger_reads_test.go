package scripts_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Count entries to the actual backing methods without modifying runtime code.
func TestDebuggerBufferedReadCounts(t *testing.T) {
	lookup := methodEntry(t, "pkg/runtime/array.go", "LookupAt")
	at := methodEntry(t, "pkg/runtime/array.go", "At")
	dataset := methodEntry(t, "pkg/vm/internal/data/dataset.go", "BufferedAt")
	for _, tc := range []struct {
		path                     string
		arrayReads, datasetReads int
	}{
		{"array/describe", 0, 0}, {"array/tail", 20, 0},
		{"dataset/describe", 0, 0}, {"dataset/tail", 20, 20},
	} {
		t.Run(tc.path, func(t *testing.T) {
			profile := filepath.Join(t.TempDir(), "reads.cover")
			cmd := exec.CommandContext(t.Context(), "go", "test", "-count=1", "-covermode=count",
				"-coverpkg=./pkg/runtime,./pkg/vm/internal/data", "-coverprofile="+profile,
				"-run=^TestDebugMillionElementBufferedPage$/"+tc.path+"$", "./pkg/vm")
			cmd.Dir = ".."
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("profiling %s: %v\n%s", tc.path, err, output)
			}

			raw, err := os.ReadFile(profile)
			if err != nil {
				t.Fatal(err)
			}

			counts := string(raw)
			if got := countEntries(t, counts, lookup) + countEntries(t, counts, at); got != tc.arrayReads {
				t.Errorf("native array reads: got %d, want %d", got, tc.arrayReads)
			}
			if got := countEntries(t, counts, dataset); got != tc.datasetReads {
				t.Errorf("buffered dataset reads: got %d, want %d", got, tc.datasetReads)
			}
		})
	}
}

func methodEntry(t *testing.T, file, method string) [2]string {
	t.Helper()
	set := token.NewFileSet()
	src, err := parser.ParseFile(set, filepath.Join("..", file), nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, decl := range src.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == method {
			// Depending on the Go version, entry counters begin at the opening
			// brace or at the first statement. Find them by source position.
			at := func(pos token.Pos) string {
				p := set.Position(pos)
				return fmt.Sprintf("/%s:%d.%d,", file, p.Line, p.Column)
			}
			return [2]string{at(fn.Body.Lbrace), at(fn.Body.List[0].Pos())}
		}
	}
	t.Fatalf("missing %s in %s", method, file)
	return [2]string{}
}

func countEntries(t *testing.T, profile string, entry [2]string) int {
	t.Helper()
	for line := range strings.SplitSeq(profile, "\n") {
		if strings.Contains(line, entry[0]) || strings.Contains(line, entry[1]) {
			fields := strings.Fields(line)
			count, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil {
				t.Fatal(err)
			}
			return count
		}
	}
	t.Fatalf("entry %v missing from coverage profile", entry)
	return 0
}
