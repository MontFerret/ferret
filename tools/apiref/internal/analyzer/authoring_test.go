package analyzer_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
)

func TestEveryStructuredStdlibCommentUsesSpecsSyntax(t *testing.T) {
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}

	err = filepath.Walk(filepath.Join(root, "pkg", "stdlib"), func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return walkErr
		}

		fileset := token.NewFileSet()
		file, parseErr := parser.ParseFile(fileset, path, nil, parser.ParseComments)
		if parseErr != nil {
			return parseErr
		}

		for _, comment := range file.Comments {
			if !hasStructuredAnnotation(comment.Text()) {
				continue
			}

			if _, parseErr := api.ParseDocumentation(strings.TrimSpace(comment.Text())); parseErr != nil {
				position := fileset.Position(comment.Pos())
				t.Errorf("%s:%d: malformed Specs documentation: %v", position.Filename, position.Line, parseErr)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func hasStructuredAnnotation(documentation string) bool {
	for _, annotation := range []string{"@param", "@return", "@throws", "@deprecated"} {
		if strings.Contains(documentation, annotation) {
			return true
		}
	}

	return false
}
