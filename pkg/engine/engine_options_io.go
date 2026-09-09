package engine

import (
	"fmt"
	"io"
	"strings"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

// WithFSRoot sets the root directory for the engine's file system.
func WithFSRoot(root string) Option {
	return gooptions.New(func(opts *config, root string) {
		opts.fsRoot = strings.TrimSpace(root)
	}).
		Value(root).
		Named("fs root").
		Validators(gooptions.NotBlank[string]()).
		Build()
}

// WithFSReadOnly sets the engine's file system to read-only mode.
func WithFSReadOnly() Option {
	return gooptions.New(func(opts *config, readOnly bool) {
		opts.fsReadOnly = readOnly
	}).
		Value(true).
		Build()
}

// WithNetwork sets the engine network service used by derived executions.
// If a network is provided, the engine will use it directly and will not manage its lifecycle.
// The host application is responsible for closing the network when it is no longer needed.
func WithNetwork(network ferretnet.Network) Option {
	return func(opts *config) error {
		if network == nil {
			return fmt.Errorf("network cannot be nil")
		}

		opts.network = network

		return opts.resources.Borrow(resource.Network)
	}
}

// WithNetworkOptions creates an Option that constructs a new network service using the provided Ferret network options.
// The engine owns the resulting service and cleans it up when replaced, on
// construction failure, or at shutdown. If no options are provided, it is a no-op.
func WithNetworkOptions(setters ...ferretnet.Option) Option {
	return func(opts *config) error {
		if len(setters) == 0 {
			return nil
		}

		net, err := ferretnet.New(setters...)
		if err != nil {
			return fmt.Errorf("create network: %w", err)
		}

		opts.network = net

		return opts.resources.Own(resource.Network, func() error {
			ferretnet.CloseIdleNetworkConnections(net)

			return nil
		})
	}
}

// WithLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithLog(writer io.Writer) Option {
	return func(opts *config) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.logger = append(opts.logger, logging.WithWriter(writer))

		return nil
	}
}

// WithLogLevel sets the logging level for the engine.
// The logging level determines the severity of log messages that will be recorded.
func WithLogLevel(lvl logging.LogLevel) Option {
	return func(opts *config) error {
		if lvl < logging.TraceLevel || lvl > logging.Disabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.logger = append(opts.logger, logging.WithLevel(lvl))

		return nil
	}
}

// WithLogFields sets the fields to be included in log entries.
// These fields can provide additional context for debugging and monitoring purposes.
func WithLogFields(fields map[string]any) Option {
	return func(opts *config) error {
		if len(fields) == 0 {
			return nil
		}

		opts.logger = append(opts.logger, logging.WithFields(fields))

		return nil
	}
}

// WithEncodingRegistry sets a custom encoding registry for query execution.
func WithEncodingRegistry(registry *encoding.Registry) Option {
	return gooptions.New(func(opts *config, registry *encoding.Registry) {
		opts.encoding = registry
	}).
		Value(registry).
		Named("encoding registry").
		Validators(gooptions.NotNilPtr[encoding.Registry]()).
		Build()
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return func(opts *config) error {
		if codec == nil {
			return encoding.ErrNilCodec
		}

		if opts.encoding == nil {
			opts.encoding = encoding.NewRegistry()
		}

		return opts.encoding.Register(encodingCodecAlias{
			Codec:       codec,
			contentType: contentType,
		})
	}
}
