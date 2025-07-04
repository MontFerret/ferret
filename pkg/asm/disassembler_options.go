package asm

type disassemblerOptions struct {
	debug bool
}

type DisassemblerOption func(*disassemblerOptions)

func newDisassemblerOptions(opts ...DisassemblerOption) *disassemblerOptions {
	options := &disassemblerOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithDebug() DisassemblerOption {
	return func(opts *disassemblerOptions) {
		opts.debug = true
	}
}
