package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func resolveExecutablePath() string {
	var res string

	switch runtime.GOOS {
	case "darwin":
		for _, c := range []string{
			"/Applications/Google Chrome Canary.app",
			"/Applications/Google Chrome.app",
		} {
			// MacOS apps are actually folders
			if info, err := os.Stat(c); err == nil && info.IsDir() {
				res = fmt.Sprintf("open %s -n", c)
				break
			}
		}

	case "linux":
		for _, c := range []string{
			"headless_shell",
			"chromium",
			"google-chrome-beta",
			"google-chrome-unstable",
			"google-chrome-stable"} {
			if _, err := exec.LookPath(c); err == nil {
				res = c
				break
			}
		}

	case "windows":
	}

	return res
}
