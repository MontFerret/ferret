package internal

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/antlr4-go/antlr/v4"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	Ctx *CompilerContext
}

func NewVisitor(src *file.Source, errors *diagnostics.ErrorHandler) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Ctx = NewCompilerContext(src, errors)

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	v.Ctx.StmtCompiler.Compile(ctx.Body())

	return nil
}

func (v *Visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func TestTypeTrackingIntegration(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		variableName string
		expectedType core.ValueType
	}{
		{"String variable", `LET x = "hello" RETURN x`, "x", core.TypeString},
		{"Integer variable", `LET x = 42 RETURN x`, "x", core.TypeInt},
		{"Float variable", `LET x = 3.14 RETURN x`, "x", core.TypeFloat},
		{"Boolean variable", `LET x = TRUE RETURN x`, "x", core.TypeBool},
		{"Array variable", `LET x = [1, 2, 3] RETURN x`, "x", core.TypeList},
		{"Object variable", `LET x = {name: "test"} RETURN x`, "x", core.TypeMap},
		{"Function call variable", `LET x = TYPENAME(1) RETURN x`, "x", core.TypeUnknown},
		{"Expression variable", `LET x = 1 + 2 RETURN x`, "x", core.TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a compiler context to test the type tracking
			src := file.NewSource("test", tt.query)
			errorHandler := diagnostics.NewErrorHandler(src, 10)

			// Create lexer and parser
			input := antlr.NewInputStream(tt.query)
			lexer := fql.NewFqlLexer(input)
			stream := antlr.NewCommonTokenStream(lexer, 0)
			parser := fql.NewFqlParser(stream)

			// Parse the query
			tree := parser.Program()

			// Create a visitor and compile
			visitor := NewVisitor(src, errorHandler)
			visitor.VisitProgram(tree.(*fql.ProgramContext))

			// Check if there were any compilation errors
			if errorHandler.HasErrors() {
				t.Fatalf("Compilation failed with errors")
			}

			// Check the debug view to see if the variable type was tracked
			debugView := visitor.Ctx.Symbols.DebugView()
			found := false
			for _, line := range debugView {
				if strings.Contains(line, tt.variableName) {
					found = true
					t.Logf("Variable %s debug info: %s", tt.variableName, line)
					
					// Look up the variable directly to check its type
					if variable, exists := visitor.Ctx.Symbols.Lookup(tt.variableName); exists {
						if variable.Type == tt.expectedType {
							t.Logf("Type tracking successful for %s: expected %v, got %v", tt.variableName, tt.expectedType, variable.Type)
							return
						}
						t.Errorf("Variable %s has wrong type: expected %v, got %v", tt.variableName, tt.expectedType, variable.Type)
					} else {
						t.Errorf("Variable %s not found in lookup", tt.variableName)
					}
					return
				}
			}

			if !found {
				t.Errorf("Variable %s not found in symbol table", tt.variableName)
			}
		})
	}
}