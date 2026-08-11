package analyzer

import (
	"fmt"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
)

func parseSignatureDocumentation(registered registeredSignature, declaration *sourceDeclaration) (api.Signature, error) {
	if strings.TrimSpace(declaration.Documentation) == "" {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "Go declaration", "has no Ferret API documentation")
	}

	documentation, err := api.ParseDocumentation(declaration.Documentation)
	if err != nil {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "Go declaration", fmt.Sprintf("has malformed Ferret API documentation: %v", err))
	}

	if strings.TrimSpace(documentation.Description) == "" {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "Go declaration", "must document non-empty Ferret-facing prose")
	}

	if documentation.Return == nil {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "Go declaration", "must document exactly one @return")
	}

	if registered.Variadic {
		if len(documentation.Parameters) == 0 {
			return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "Go declaration", "variadic registration must document at least one logical parameter")
		}
	} else if len(documentation.Parameters) != registered.Arity {
		return api.Signature{}, diagnostic(
			declaration,
			registered.QualifiedName,
			"Go declaration",
			fmt.Sprintf("documents %d parameters but runtime registration has fixed arity %d", len(documentation.Parameters), registered.Arity),
		)
	}

	return api.Signature{
		Parameters:  copyParameters(documentation.Parameters),
		Variadic:    registered.Variadic,
		Description: documentation.Description,
		Return:      documentation.Return,
		Throws:      append([]api.Throw(nil), documentation.Throws...),
		Deprecated:  documentation.Deprecated,
	}, nil
}

func parseAssertionDocumentation(registered registeredSignature, descriptor assertionDescriptor) (api.Signature, error) {
	declaration := descriptor.Declaration
	if strings.TrimSpace(declaration.Documentation) == "" {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "assertion descriptor", "has no Ferret API documentation")
	}

	documentation, err := api.ParseDocumentation(declaration.Documentation)
	if err != nil {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "assertion descriptor", fmt.Sprintf("has malformed Ferret API documentation: %v", err))
	}

	if strings.TrimSpace(documentation.Description) == "" {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "assertion descriptor", "must document namespace-neutral Ferret-facing prose")
	}

	if documentation.Return == nil {
		return api.Signature{}, diagnostic(declaration, registered.QualifiedName, "assertion descriptor", "must document exactly one @return")
	}

	if len(documentation.Parameters) != descriptor.Max {
		return api.Signature{}, diagnostic(
			declaration,
			registered.QualifiedName,
			"assertion descriptor",
			fmt.Sprintf("documents %d parameters but Args.Max is %d", len(documentation.Parameters), descriptor.Max),
		)
	}

	if registered.Variadic || registered.Arity < descriptor.Min || registered.Arity > descriptor.Max {
		return api.Signature{}, diagnostic(
			declaration,
			registered.QualifiedName,
			"assertion descriptor",
			fmt.Sprintf("runtime signature does not match literal Args range %d..%d", descriptor.Min, descriptor.Max),
		)
	}

	return api.Signature{
		Parameters:  copyParameters(documentation.Parameters[:registered.Arity]),
		Description: documentation.Description,
		Return:      documentation.Return,
		Throws:      append([]api.Throw(nil), documentation.Throws...),
		Deprecated:  documentation.Deprecated,
	}, nil
}

func copyParameters(parameters []api.Parameter) []api.Parameter {
	result := make([]api.Parameter, len(parameters))
	copy(result, parameters)

	return result
}

func diagnostic(declaration *sourceDeclaration, qualifiedName, subject, reason string) error {
	return fmt.Errorf(
		"%s:%d: package %s; Ferret %s; %s %s: %s",
		declaration.File,
		declaration.Line,
		declaration.PackagePath,
		qualifiedName,
		subject,
		declaration.Name,
		reason,
	)
}
