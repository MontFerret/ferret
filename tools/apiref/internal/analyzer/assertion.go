package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

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
		if !ok || identifierName(field.Key) != "Args" {
			continue
		}

		arguments, ok := field.Value.(*ast.CompositeLit)
		if !ok {
			return 0, 0, fmt.Errorf("has non-literal Args bounds")
		}

		minimum, minimumSet := 0, false
		maximum, maximumSet := 0, false

		for _, argumentElement := range arguments.Elts {
			argument, ok := argumentElement.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			switch identifierName(argument.Key) {
			case "Min":
				value, err := integerLiteral(argument.Value)
				if err != nil {
					return 0, 0, fmt.Errorf("has invalid Args.Min: %w", err)
				}

				minimum, minimumSet = value, true
			case "Max":
				value, err := integerLiteral(argument.Value)
				if err != nil {
					return 0, 0, fmt.Errorf("has invalid Args.Max: %w", err)
				}

				maximum, maximumSet = value, true
			}
		}

		if !minimumSet || !maximumSet {
			return 0, 0, fmt.Errorf("must declare literal Args.Min and Args.Max")
		}

		if minimum < 0 || maximum < minimum || maximum > 4 {
			return 0, 0, fmt.Errorf("declares unsupported Args range %d..%d", minimum, maximum)
		}

		return minimum, maximum, nil
	}

	return 0, 0, fmt.Errorf("does not declare Args bounds")
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
