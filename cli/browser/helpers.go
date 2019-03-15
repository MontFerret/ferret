package browser

import (
	"os"
	"os/exec"
	"runtime"
)

func resolveExecutablePath() (path string) {

	switch runtime.GOOS {
	case "darwin":
		for _, c := range []string{
			"/Applications/Google Chrome Canary.app",
			"/Applications/Google Chrome.app",
		} {
			// MacOS apps are actually folders
			info, err := os.Stat(c)
			if err == nil && info.IsDir() {
				path = c
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
				path = c
				break
			}
		}

	case "windows":
	}

	return
}
