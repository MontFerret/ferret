package analyzer

import (
	"fmt"
	"reflect"
	goruntime "runtime"
	"sort"
	"strings"

	ferretruntime "github.com/MontFerret/ferret/v2/pkg/runtime"
)

func registeredSignatures(functions *ferretruntime.Functions) ([]registeredSignature, error) {
	if functions == nil {
		return nil, fmt.Errorf("stdlib function registry is nil")
	}

	type registrationKey struct {
		name     string
		arity    int
		variadic bool
	}

	signatures := make([]registeredSignature, 0, functions.Size())
	seen := make(map[registrationKey]struct{}, functions.Size())
	var registrationErr error
	appendSignature := func(name string, arity int, variadic bool, fn any) bool {
		if registrationErr != nil {
			return false
		}

		if err := validateRegisteredName(name); err != nil {
			registrationErr = err

			return false
		}

		key := registrationKey{name: name, arity: arity, variadic: variadic}
		if _, exists := seen[key]; exists {
			registrationErr = fmt.Errorf("duplicate runtime function registration %s/%d", name, arity)

			return false
		}

		seen[key] = struct{}{}

		signature, err := registeredSignatureFor(name, arity, variadic, fn)
		if err != nil {
			registration := fmt.Sprintf("%s/%d", name, arity)
			if variadic {
				registration = name + "/variadic"
			}

			registrationErr = fmt.Errorf("resolve runtime function %s: %w", registration, err)

			return false
		}

		signatures = append(signatures, signature)

		return true
	}

	functions.A0().ForEach(func(fn ferretruntime.Function0, name string) bool {
		return appendSignature(name, 0, false, fn)
	})
	functions.A1().ForEach(func(fn ferretruntime.Function1, name string) bool {
		return appendSignature(name, 1, false, fn)
	})
	functions.A2().ForEach(func(fn ferretruntime.Function2, name string) bool {
		return appendSignature(name, 2, false, fn)
	})
	functions.A3().ForEach(func(fn ferretruntime.Function3, name string) bool {
		return appendSignature(name, 3, false, fn)
	})
	functions.A4().ForEach(func(fn ferretruntime.Function4, name string) bool {
		return appendSignature(name, 4, false, fn)
	})
	functions.Var().ForEach(func(fn ferretruntime.Function, name string) bool {
		return appendSignature(name, 0, true, fn)
	})

	if registrationErr != nil {
		return nil, registrationErr
	}

	if len(signatures) != functions.Size() {
		return nil, fmt.Errorf("runtime registry reported %d signatures but enumerated %d", functions.Size(), len(signatures))
	}

	sort.Slice(signatures, func(i, j int) bool {
		if signatures[i].QualifiedName != signatures[j].QualifiedName {
			return signatures[i].QualifiedName < signatures[j].QualifiedName
		}

		if signatures[i].Variadic != signatures[j].Variadic {
			return !signatures[i].Variadic
		}

		return signatures[i].Arity < signatures[j].Arity
	})

	return signatures, nil
}

func validateRegisteredName(name string) error {
	_, terminal := splitQualifiedName(name)
	if terminal == "" {
		return fmt.Errorf("runtime function %q has an empty terminal name", name)
	}

	for i := 0; i < len(terminal); i++ {
		if terminal[i] >= 'A' && terminal[i] <= 'Z' {
			return fmt.Errorf("runtime function %q terminal name must be canonical lowercase", name)
		}
	}

	return nil
}

func registeredSignatureFor(name string, arity int, variadic bool, fn any) (registeredSignature, error) {
	value := reflect.ValueOf(fn)
	if !value.IsValid() || value.Kind() != reflect.Func || value.IsNil() {
		return registeredSignature{}, fmt.Errorf("runtime function is nil or is not a function")
	}

	pointer := value.Pointer()
	resolved := goruntime.FuncForPC(pointer)
	if resolved == nil {
		return registeredSignature{}, fmt.Errorf("runtime function has no source symbol")
	}

	symbol := strings.TrimSuffix(resolved.Name(), "-fm")
	file, line := resolved.FileLine(pointer)
	namespace, publicName := splitQualifiedName(name)

	return registeredSignature{
		QualifiedName: name,
		Namespace:     namespace,
		Name:          publicName,
		Symbol:        symbol,
		PackagePath:   symbolPackage(symbol),
		File:          file,
		Line:          line,
		Arity:         arity,
		Variadic:      variadic,
	}, nil
}

func symbolPackage(symbol string) string {
	packageStart := strings.LastIndex(symbol, "/") + 1
	separator := strings.Index(symbol[packageStart:], ".")
	if separator < 0 {
		return ""
	}

	return symbol[:packageStart+separator]
}

func splitQualifiedName(name string) (string, string) {
	index := strings.LastIndex(name, ferretruntime.NamespaceSeparator)
	if index < 0 {
		return "", name
	}

	return name[:index], name[index+len(ferretruntime.NamespaceSeparator):]
}
