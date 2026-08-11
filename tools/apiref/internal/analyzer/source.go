package analyzer

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

func loadSourceCatalog(ctx context.Context, root string) (*sourceCatalog, error) {
	configuration := &packages.Config{
		Context: ctx,
		Dir:     root,
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule,
		Tests: false,
	}

	loaded, err := packages.Load(configuration, "./pkg/stdlib/...")
	if err != nil {
		return nil, fmt.Errorf("load pkg/stdlib source: %w", err)
	}

	catalog := &sourceCatalog{
		Declarations: make(map[string]*sourceDeclaration),
		Assertions:   make(map[string]assertionDescriptor),
	}

	for _, pkg := range loaded {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("load package %s: %s", pkg.PkgPath, packageErrors(pkg.Errors))
		}

		for _, file := range pkg.Syntax {
			if err := collectDeclarations(catalog, pkg, file); err != nil {
				return nil, err
			}
		}
	}

	for _, pkg := range loaded {
		if !strings.HasSuffix(pkg.PkgPath, "/pkg/stdlib/testing") {
			continue
		}

		for _, file := range pkg.Syntax {
			if err := collectAssertionRegistrations(catalog, pkg, file); err != nil {
				return nil, err
			}
		}
	}

	return catalog, nil
}

func collectDeclarations(catalog *sourceCatalog, pkg *packages.Package, file *ast.File) error {
	for _, declaration := range file.Decls {
		switch declaration := declaration.(type) {
		case *ast.FuncDecl:
			if declaration.Recv != nil {
				continue
			}

			object, ok := pkg.TypesInfo.Defs[declaration.Name].(*types.Func)
			if !ok || object.Pkg() == nil {
				continue
			}

			entry := sourceDeclarationFor(pkg, declaration.Name, declaration.Doc, returnedExpression(declaration.Body))
			entry.Key = object.Pkg().Path() + "." + object.Name()

			if err := addSourceDeclaration(catalog, entry); err != nil {
				return err
			}
		case *ast.GenDecl:
			if declaration.Tok != token.VAR {
				continue
			}

			for _, specification := range declaration.Specs {
				value, ok := specification.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for index, name := range value.Names {
					object, ok := pkg.TypesInfo.Defs[name].(*types.Var)
					if !ok || object.Pkg() == nil {
						continue
					}

					comment := value.Doc
					if comment == nil {
						comment = declaration.Doc
					}

					var expression ast.Expr
					if index < len(value.Values) {
						expression = value.Values[index]
					} else if len(value.Values) == 1 {
						expression = value.Values[0]
					}

					entry := sourceDeclarationFor(pkg, name, comment, expression)
					entry.Key = object.Pkg().Path() + "." + object.Name()
					entry.Factory = calledFunctionKey(pkg, expression)

					if err := addSourceDeclaration(catalog, entry); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func sourceDeclarationFor(pkg *packages.Package, name *ast.Ident, comment *ast.CommentGroup, expression ast.Expr) *sourceDeclaration {
	position := pkg.Fset.Position(name.Pos())
	documentation := ""
	if comment != nil {
		documentation = strings.TrimSpace(comment.Text())
	}

	file := position.Filename
	if pkg.Module != nil && pkg.Module.Dir != "" {
		if relative, err := filepath.Rel(pkg.Module.Dir, file); err == nil {
			file = filepath.ToSlash(relative)
		}
	}

	return &sourceDeclaration{
		Name:          name.Name,
		PackagePath:   pkg.PkgPath,
		File:          file,
		Line:          position.Line,
		Documentation: documentation,
		Expression:    expression,
	}
}

func addSourceDeclaration(catalog *sourceCatalog, declaration *sourceDeclaration) error {
	if previous, exists := catalog.Declarations[declaration.Key]; exists {
		return fmt.Errorf("ambiguous Go declaration %s at %s:%d and %s:%d", declaration.Key, previous.File, previous.Line, declaration.File, declaration.Line)
	}

	catalog.Declarations[declaration.Key] = declaration

	return nil
}

func collectAssertionRegistrations(catalog *sourceCatalog, pkg *packages.Package, file *ast.File) error {
	var collectionErr error

	ast.Inspect(file, func(node ast.Node) bool {
		if collectionErr != nil {
			return false
		}

		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		function, ok := call.Fun.(*ast.Ident)
		if !ok || function.Name != "registerPositive" && function.Name != "registerNegative" {
			return true
		}

		if len(call.Args) != 3 {
			collectionErr = fmt.Errorf("%s: assertion registration %s must have exactly three arguments", pkg.Fset.Position(call.Pos()), function.Name)

			return false
		}

		nameLiteral, ok := call.Args[1].(*ast.BasicLit)
		if !ok || nameLiteral.Kind != token.STRING {
			collectionErr = fmt.Errorf("%s: assertion registration %s must use a string literal public name", pkg.Fset.Position(call.Pos()), function.Name)

			return false
		}

		name, err := strconv.Unquote(nameLiteral.Value)
		if err != nil {
			collectionErr = fmt.Errorf("%s: decode assertion public name: %w", pkg.Fset.Position(nameLiteral.Pos()), err)

			return false
		}

		descriptorName, ok := call.Args[2].(*ast.Ident)
		if !ok {
			collectionErr = fmt.Errorf("%s: assertion %s must use a statically named descriptor", pkg.Fset.Position(call.Args[2].Pos()), name)

			return false
		}

		object, ok := pkg.TypesInfo.Uses[descriptorName].(*types.Var)
		if !ok || object.Pkg() == nil {
			collectionErr = fmt.Errorf("%s: resolve assertion descriptor %s", pkg.Fset.Position(descriptorName.Pos()), descriptorName.Name)

			return false
		}

		declarationKey := object.Pkg().Path() + "." + object.Name()
		declaration, exists := catalog.Declarations[declarationKey]
		if !exists {
			collectionErr = fmt.Errorf("%s: declaration for assertion descriptor %s was not loaded", pkg.Fset.Position(descriptorName.Pos()), declarationKey)

			return false
		}

		minimum, maximum, err := assertionBounds(catalog, declaration)
		if err != nil {
			collectionErr = diagnostic(declaration, assertionQualifiedName(function.Name, name), "assertion descriptor", err.Error())

			return false
		}

		qualifiedName := assertionQualifiedName(function.Name, name)
		if _, exists := catalog.Assertions[qualifiedName]; exists {
			collectionErr = diagnostic(declaration, qualifiedName, "assertion descriptor", "is registered more than once")

			return false
		}

		catalog.Assertions[qualifiedName] = assertionDescriptor{
			Declaration: declaration,
			Min:         minimum,
			Max:         maximum,
		}

		return true
	})

	return collectionErr
}

func returnedExpression(body *ast.BlockStmt) ast.Expr {
	if body == nil {
		return nil
	}

	for _, statement := range body.List {
		returnStatement, ok := statement.(*ast.ReturnStmt)
		if ok && len(returnStatement.Results) == 1 {
			return returnStatement.Results[0]
		}
	}

	return nil
}

func calledFunctionKey(pkg *packages.Package, expression ast.Expr) string {
	call, ok := expression.(*ast.CallExpr)
	if !ok {
		return ""
	}

	var identifier *ast.Ident
	switch function := call.Fun.(type) {
	case *ast.Ident:
		identifier = function
	case *ast.SelectorExpr:
		identifier = function.Sel
	}

	if identifier == nil {
		return ""
	}

	object, ok := pkg.TypesInfo.Uses[identifier].(*types.Func)
	if !ok || object.Pkg() == nil {
		return ""
	}

	return object.Pkg().Path() + "." + object.Name()
}

func assertionQualifiedName(registration, name string) string {
	if registration == "registerNegative" {
		return "T::NOT::" + name
	}

	return "T::" + name
}

func packageErrors(errors []packages.Error) string {
	messages := make([]string, 0, len(errors))
	for _, err := range errors {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, "; ")
}
