package browser

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

func Launch(setters ...Option) (*Browser, error) {
	opts := &Options{
		headless:         false,
		debuggingAddress: "0.0.0.0",
		debuggingPort:    9222,
	}

	for _, setter := range setters {
		setter(opts)
	}

	var flags Flags

	if !opts.ignoreDefaultArgs {
		flags = DefaultFlags()
	} else {
		flags = Flags{}
	}

	flags.Set("remote-debugging-port", opts.debuggingPort)

	if opts.devtools {
		flags.SetN("auto-open-devtools-for-tabs")
	}

	if opts.headless {
		flags.SetN("headless")
		flags.SetN("hide-scrollbars")
		flags.SetN("mute-audio")
	}

	if runtime.GOOS == "windows" {
		flags.SetN("disable-gpu")
	}

	temporaryUserDataDir := opts.userDataDir

	if temporaryUserDataDir == "" && opts.noUserDataDir == false {
		dirName, err := ioutil.TempDir(os.TempDir(), "ferret_dev_profile-")

		if err != nil {
			return nil, err
		}

		temporaryUserDataDir = dirName
	}

	workDir := filepath.Join(os.TempDir(), "ferret-chrome")

	err := os.MkdirAll(workDir, 0700)

	if err != nil {
		return nil, errors.Errorf("cannot create working directory '%s'", workDir)
	}

	if temporaryUserDataDir != "" {
		flags.Set("user-data-dir", temporaryUserDataDir)
	}

	chromeExecutable := opts.executablePath

	if chromeExecutable == "" {
		chromeExecutable = resolveExecutablePath()

		if chromeExecutable == "" {
			return nil, errors.New("Chrome not found")
		}
	}

	execArgs := []string{chromeExecutable, "--args"}
	execArgs = append(execArgs, flags.List()...)

	cmd := exec.Command("open", execArgs...)
	cmd.Dir = workDir

	err = cmd.Start()

	if err != nil {
		return nil, err
	}

	if opts.dumpio {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return &Browser{cmd, flags}, nil
}

func DefaultFlags() Flags {
	return Flags{
		"disable-background-networking":          nil,
		"disable-background-timer-throttling":    nil,
		"disable-breakpad":                       nil,
		"disable-client-side-phishing-detection": nil,
		"disable-default-apps":                   nil,
		"disable-dev-shm-usage":                  nil,
		"disable-extensions":                     nil,
		"disable-features=site-per-process":      nil,
		"disable-hang-monitor":                   nil,
		"disable-popup-blocking":                 nil,
		"disable-prompt-on-repost":               nil,
		"disable-sync":                           nil,
		"disable-translate":                      nil,
		"metrics-recording-only":                 nil,
		"no-first-run":                           nil,
		"safebrowsing-disable-auto-update":       nil,
		"enable-automation":                      nil,
		"password-store=basic":                   nil,
		"use-mock-keychain":                      nil,
	}
}
