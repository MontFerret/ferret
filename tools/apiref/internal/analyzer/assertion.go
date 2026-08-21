package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/packages"
)

func collectAssertionCatalog(catalog *sourceCatalog, pkg *packages.Package, file *ast.File) (bool, error) {
	for _, declaration := range file.Decls {
		variables, ok := declaration.(*ast.GenDecl)
		if !ok || variables.Tok != token.VAR {
			continue
		}

		for _, specification := range variables.Specs {
			value, ok := specification.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range value.Names {
				if name.Name != "assertionCatalog" {
					continue
				}

				if len(value.Names) != 1 || len(value.Values) != 1 {
					return true, fmt.Errorf("%s: assertionCatalog must be declared as one literal value", pkg.Fset.Position(value.Pos()))
				}

				entries, ok := value.Values[0].(*ast.CompositeLit)
				if !ok || !isAssertionCatalogType(entries.Type) {
					return true, fmt.Errorf("%s: assertionCatalog must be a []assertionRegistration literal", pkg.Fset.Position(value.Values[0].Pos()))
				}

				if len(entries.Elts) == 0 {
					return true, fmt.Errorf("%s: assertionCatalog must not be empty", pkg.Fset.Position(entries.Pos()))
				}

				for _, element := range entries.Elts {
					entry, ok := element.(*ast.CompositeLit)
					if !ok {
						return true, fmt.Errorf("%s: assertionCatalog entries must be keyed literals", pkg.Fset.Position(element.Pos()))
					}

					if err := collectAssertionCatalogEntry(catalog, pkg, entry); err != nil {
						return true, err
					}
				}

				return true, nil
			}
		}
	}

	return false, nil
}

func collectAssertionCatalogEntry(catalog *sourceCatalog, pkg *packages.Package, entry *ast.CompositeLit) error {
	fields := make(map[string]ast.Expr, 3)
	for _, element := range entry.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			return fmt.Errorf("%s: assertionCatalog entries must use keyed fields", pkg.Fset.Position(element.Pos()))
		}

		name := identifierName(field.Key)
		switch name {
		case "name", "descriptor", "negatable":
		default:
			return fmt.Errorf("%s: assertionCatalog entry has unsupported field %q", pkg.Fset.Position(field.Key.Pos()), name)
		}

		if _, exists := fields[name]; exists {
			return fmt.Errorf("%s: assertionCatalog entry repeats field %q", pkg.Fset.Position(field.Key.Pos()), name)
		}

		fields[name] = field.Value
	}

	for _, name := range []string{"name", "descriptor", "negatable"} {
		if _, exists := fields[name]; !exists {
			return fmt.Errorf("%s: assertionCatalog entry must declare %s", pkg.Fset.Position(entry.Pos()), name)
		}
	}

	nameLiteral, ok := fields["name"].(*ast.BasicLit)
	if !ok || nameLiteral.Kind != token.STRING {
		return fmt.Errorf("%s: assertionCatalog name must be a string literal", pkg.Fset.Position(fields["name"].Pos()))
	}

	name, err := strconv.Unquote(nameLiteral.Value)
	if err != nil {
		return fmt.Errorf("%s: decode assertionCatalog name: %w", pkg.Fset.Position(nameLiteral.Pos()), err)
	}

	if name == "" {
		return fmt.Errorf("%s: assertionCatalog name must not be empty", pkg.Fset.Position(nameLiteral.Pos()))
	}

	descriptorName, ok := fields["descriptor"].(*ast.Ident)
	if !ok {
		return fmt.Errorf("%s: assertion %s must use a statically named descriptor", pkg.Fset.Position(fields["descriptor"].Pos()), name)
	}

	object, ok := pkg.TypesInfo.Uses[descriptorName].(*types.Var)
	if !ok || object.Pkg() == nil {
		return fmt.Errorf("%s: resolve assertion descriptor %s", pkg.Fset.Position(descriptorName.Pos()), descriptorName.Name)
	}

	declarationKey := object.Pkg().Path() + "." + object.Name()
	declaration, exists := catalog.Declarations[declarationKey]
	if !exists {
		return fmt.Errorf("%s: declaration for assertion descriptor %s was not loaded", pkg.Fset.Position(descriptorName.Pos()), declarationKey)
	}

	minimum, maximum, err := assertionBounds(catalog, declaration)
	if err != nil {
		return diagnostic(declaration, "t::"+name, "assertion descriptor", err.Error())
	}

	negatableName, ok := fields["negatable"].(*ast.Ident)
	if !ok || negatableName.Name != "true" && negatableName.Name != "false" {
		return fmt.Errorf("%s: assertionCatalog negatable must be a boolean literal", pkg.Fset.Position(fields["negatable"].Pos()))
	}

	descriptor := assertionDescriptor{
		Declaration: declaration,
		Min:         minimum,
		Max:         maximum,
	}
	if err := addCatalogAssertion(catalog, "t::"+name, descriptor); err != nil {
		return err
	}

	if negatableName.Name == "true" {
		if err := addCatalogAssertion(catalog, "t::not::"+name, descriptor); err != nil {
			return err
		}
	}

	return nil
}

func addCatalogAssertion(catalog *sourceCatalog, qualifiedName string, descriptor assertionDescriptor) error {
	if _, exists := catalog.Assertions[qualifiedName]; exists {
		return diagnostic(descriptor.Declaration, qualifiedName, "assertion descriptor", "is registered more than once")
	}

	catalog.Assertions[qualifiedName] = descriptor

	return nil
}

func isAssertionCatalogType(expression ast.Expr) bool {
	slice, ok := expression.(*ast.ArrayType)
	if !ok || slice.Len != nil {
		return false
	}

	return identifierName(slice.Elt) == "assertionRegistration"
}

func assertionBounds(catalog *sourceCatalog, declaration *sourceDeclaration) (int, int, error) {
	switch declaration.Expression.(type) {
	case *ast.CallExpr:
		if declaration.Factory == "" {
			return 0, 0, fmt.Errorf("uses an unresolved assertion factory")
		}

		factory, exists := catalog.Declarations[declaration.Factory]
		if !exists {
			return 0, 0, fmt.Errorf("uses assertion factory %s whose source declaration was not loaded", declaration.Factory)
		}

		return assertionBounds(catalog, factory)
	case *ast.CompositeLit:
		return compositeAssertionBounds(declaration.Expression.(*ast.CompositeLit))
	default:
		return 0, 0, fmt.Errorf("uses an unsupported dynamic assertion expression")
	}
}

func compositeAssertionBounds(assertion *ast.CompositeLit) (int, int, error) {
	for _, element := range assertion.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok || identifierName(field.Key) != "args" {
			continue
		}

		arguments, ok := field.Value.(*ast.CompositeLit)
		if !ok {
			return 0, 0, fmt.Errorf("has non-literal args bounds")
		}

		minimum, minimumSet := 0, false
		maximum, maximumSet := 0, false

		for _, argumentElement := range arguments.Elts {
			argument, ok := argumentElement.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			switch identifierName(argument.Key) {
			case "min":
				value, err := integerLiteral(argument.Value)
				if err != nil {
					return 0, 0, fmt.Errorf("has invalid args.min: %w", err)
				}

				minimum, minimumSet = value, true
			case "max":
				value, err := integerLiteral(argument.Value)
				if err != nil {
					return 0, 0, fmt.Errorf("has invalid args.max: %w", err)
				}

				maximum, maximumSet = value, true
			}
		}

		if !minimumSet || !maximumSet {
			return 0, 0, fmt.Errorf("must declare literal args.min and args.max")
		}

		if minimum < 0 || maximum < minimum || maximum > 4 {
			return 0, 0, fmt.Errorf("declares unsupported args range %d..%d", minimum, maximum)
		}

		return minimum, maximum, nil
	}

	return 0, 0, fmt.Errorf("does not declare args bounds")
}

func integerLiteral(expression ast.Expr) (int, error) {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != token.INT {
		return 0, fmt.Errorf("expected integer literal")
	}

	value, err := strconv.Atoi(literal.Value)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", literal.Value, err)
	}

	return value, nil
}

func identifierName(expression ast.Expr) string {
	identifier, _ := expression.(*ast.Ident)
	if identifier == nil {
		return ""
	}

	return identifier.Name
}
