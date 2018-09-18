package browser

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"runtime"
)

type Browser struct {
	cmd   *exec.Cmd
	flags Flags
}

func (b *Browser) Flags() Flags {
	return b.flags
}

func (b *Browser) DebuggingAddress() string {
	if !b.Flags().Has("remote-debugging-address") {
		b.Flags().Set("remote-debugging-address", "0.0.0.0")
	}

	value, _ := b.Flags().Get("remote-debugging-address")

	return value.(string)
}

func (b *Browser) DebuggingPort() int {
	if !b.Flags().Has("remote-debugging-port") {
		b.Flags().Set("remote-debugging-port", 9222)
	}

	value, _ := b.Flags().Get("remote-debugging-port")

	return value.(int)
}

func (b *Browser) Close() error {
	var err error

	if runtime.GOOS != "windows" {
		err = b.cmd.Process.Signal(os.Interrupt)
	} else {
		err = b.cmd.Process.Kill()
	}

	_, err = b.cmd.Process.Wait()

	if err != nil {
		return errors.Wrap(err, "error waiting for process exit, result unknown")
	}

	tmpDir, err := b.flags.GetString("user-data-dir")

	if err != nil {
		return nil
	}

	os.RemoveAll(tmpDir)

	if err != nil {
		return err
	}

	return nil
}
