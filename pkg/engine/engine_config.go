package engine

import (
	"errors"
	"fmt"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

type config struct {
	library           runtime.Library
	network           ferretnet.Network
	resources         *resource.Manager
	hooks             *host.Hooks
	encoding          *encoding.Registry
	params            runtime.Params
	programLoader     *artifact.Loader
	fsRoot            string
	stdlib            stdlib.Set
	modules           []module.Module
	logger            []logging.Option
	optimizationLevel OptimizationLevel
	maxActiveSessions int
	maxIdleVMsPerPlan int
	maxVMsPerPlan     int
	fsReadOnly        bool
}

const (
	defaultMaxActiveSessions = 0 // 0 means no limit on active sessions.
	defaultVMPoolSize        = 8
	defaultMaxVMsPerPlan     = 0 // 0 means no limit on total VMs per plan.
)

func defaultConfig() config {
	return config{
		library:           runtime.NewLibrary(),
		resources:         resource.NewManager(),
		params:            make(map[string]runtime.Value),
		encoding:          encoding.NewRegistry(encodingjson.Default, encodingmsgpack.Default),
		programLoader:     artifact.NewDefaultLoader(),
		hooks:             host.NewHooks(),
		optimizationLevel: OptimizationFull,
		maxActiveSessions: defaultMaxActiveSessions,
		maxIdleVMsPerPlan: defaultVMPoolSize,
		maxVMsPerPlan:     defaultMaxVMsPerPlan,
		stdlib:            stdlib.Full(),
	}
}

func newConfig(setters []Option) (config, error) {
	opts, err := gooptions.ApplyTo(defaultConfig(), setters...)
	if err == nil {
		if registerErr := opts.stdlib.Register(opts.library); registerErr != nil {
			err = fmt.Errorf("stdlib: %w", registerErr)
		}
	}

	if err != nil {
		if closeErr := opts.resources.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return config{}, err
	}

	return opts, nil
}
