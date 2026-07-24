package sdk

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	decodeConfig struct {
		err                   error
		onlyFields            map[string]struct{}
		requiredTypes         []runtime.Type
		disallowUnknownFields bool
		disallowNoneValues    bool
	}

	// DecodeOption configures Decode and the argument decoding helpers.
	DecodeOption func(*decodeConfig)
)

// DisallowUnknownFields rejects object keys that do not match a tagged struct field.
func DisallowUnknownFields() DecodeOption {
	return func(config *decodeConfig) {
		config.disallowUnknownFields = true
	}
}

// RequireType requires the root source value to match at least one expected
// Ferret runtime type before conversion.
//
// Passing no types, or a nil type, is an SDK authoring error.
func RequireType(expected ...runtime.Type) DecodeOption {
	copied := append([]runtime.Type(nil), expected...)
	var optionErr error

	if len(copied) == 0 {
		optionErr = runtime.Error(runtime.ErrInvalidArgument, "RequireType requires at least one type")
	}

	if optionErr == nil {
		for _, expectedType := range copied {
			if expectedType == nil {
				optionErr = runtime.Error(runtime.ErrInvalidArgument, "RequireType does not accept nil types")

				break
			}
		}
	}

	return func(config *decodeConfig) {
		if optionErr != nil {
			config.err = optionErr

			return
		}

		config.requiredTypes = copied
	}
}

// OnlyFields restricts root object decoding to the named tagged struct fields.
// Names are matched case-insensitively. The option is valid only for struct targets.
func OnlyFields(names ...string) DecodeOption {
	fields := make(map[string]struct{}, len(names))
	var optionErr error

	for _, name := range names {
		if name == "" {
			optionErr = runtime.Error(runtime.ErrInvalidArgument, "OnlyFields does not accept empty field names")

			break
		}

		fields[strings.ToLower(name)] = struct{}{}
	}

	return func(config *decodeConfig) {
		if optionErr != nil {
			config.err = optionErr

			return
		}

		config.onlyFields = fields
	}
}

// DisallowNoneValues rejects explicit None values decoded into native Go values.
// Fields whose destination type is exactly runtime.Value preserve runtime.None.
func DisallowNoneValues() DecodeOption {
	return func(config *decodeConfig) {
		config.disallowNoneValues = true
	}
}
