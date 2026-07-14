package sdk

type decodeConfig struct {
	disallowUnknownFields bool
}

// DecodeOption configures Decode and the argument decoding helpers.
type DecodeOption func(*decodeConfig)

// DisallowUnknownFields rejects object keys that do not match a tagged struct field.
func DisallowUnknownFields() DecodeOption {
	return func(config *decodeConfig) {
		config.disallowUnknownFields = true
	}
}
