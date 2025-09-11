package runtime

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

// ErrorArg creates an error for an invalid argument at the specified position.
// The position is 0-based internally but reported as 1-based to users.
func ErrorArg(err error, pos int) error {
	return Errorf(
		ErrInvalidArgumentType,
		"expected argument %d to be: %s",
		pos+1, err.Error(),
	)
}

// ValidateArgs validates that the number of arguments is within the specified range.
// It returns an error if the argument count is outside the [minimum, maximum] range.
func ValidateArgs(args []Value, minimum, maximum int) error {
	count := len(args)

	if count < minimum || count > maximum {
		return Error(
			ErrInvalidArgumentNumber,
			fmt.Sprintf(
				"expected number of arguments %d-%d, but got %d",
				minimum,
				maximum,
				len(args)))
	}

	return nil
}

// ValidateArgType validates that the argument at the specified position matches
// the given type assertion. If the position is beyond the arguments array,
// no validation is performed (returns nil).
func ValidateArgType(args []Value, pos int, assertion TypeAssertion) error {
	if pos >= len(args) {
		return nil
	}

	arg := args[pos]

	err := assertion(arg)

	if err == nil {
		return nil
	}

	return ErrorArg(err, pos)
}

func makeFunctionName(namespace, name string) string {
	name = strings.ToUpper(name)

	if namespace == emptyNS {
		return name
	}

	return namespace + NamespaceSeparator + name
}

func functionsHash(f *functionRegistry) uint64 {
	if f == nil {
		return 0
	}

	names := f.Names()
	sort.Strings(names)

	hasher := fnv.New64a()

	for _, name := range names {
		_, _ = hasher.Write([]byte(name))
	}

	return hasher.Sum64()
}
