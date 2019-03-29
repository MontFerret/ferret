package browser

type (
	Option func(opts *Options)

	Options struct {
		debuggingPort     int
		debuggingAddress  string
		ignoreDefaultArgs bool
		executablePath    string
		ignoreHTTPSErrors bool
		slowMo            bool
		dumpio            bool
		headless          bool
		devtools          bool
		userDataDir       string
		noUserDataDir     bool
	}
)

const (
	goosWindows = "windows"
	goosLinux   = "linux"
	goosDarwin  = "darwin"
)

func WithoutDefaultArgs() Option {
	return func(opts *Options) {
		opts.ignoreDefaultArgs = true
	}
}

func WithCustomInstallation(executablePath string) Option {
	return func(opts *Options) {
		opts.executablePath = executablePath
	}
}

func WithIgnoredHTTPSErrors() Option {
	return func(opts *Options) {
		opts.ignoreHTTPSErrors = true
	}
}

func WithSlowMo() Option {
	return func(opts *Options) {
		opts.slowMo = true
	}
}

func WithIO() Option {
	return func(opts *Options) {
		opts.dumpio = true
	}
}

func WithHeadless() Option {
	return func(opts *Options) {
		opts.headless = true
	}
}

func WithDevtools() Option {
	return func(opts *Options) {
		opts.devtools = true
	}
}

func WithDebugginPort(num int) Option {
	return func(opts *Options) {
		opts.debuggingPort = num
	}
}

func WithUserDataDir(str string) Option {
	return func(opts *Options) {
		opts.userDataDir = str
	}
}

func WithoutUserDataDir() Option {
	return func(opts *Options) {
		opts.noUserDataDir = true
	}
}
