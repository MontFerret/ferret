package engine

import (
	"errors"

	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
)

// sessionConfig preserves the native option target while the session component
// owns the applied settings and their preparation semantics.
type sessionConfig struct {
	config enginesession.Config
}

func defaultSessionConfig() sessionConfig {
	return sessionConfig{config: enginesession.NewConfig()}
}

func newSessionConfig(setters []SessionOption) (sessionConfig, error) {
	if len(setters) == 0 {
		return defaultSessionConfig(), nil
	}

	opts := defaultSessionConfig()
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(&opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	if err := errors.Join(failures...); err != nil {
		return sessionConfig{}, err
	}

	return opts, nil
}
