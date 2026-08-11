package publisher

import (
	"fmt"
	"os"
	"path/filepath"
)

func atomicWrite(path string, data []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".api-publication-*")
	if err != nil {
		return err
	}

	temporaryPath := temporary.Name()
	removeTemporary := true
	defer func() {
		if removeTemporary {
			_ = os.Remove(temporaryPath)
		}
	}()

	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()

		return err
	}

	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()

		return err
	}

	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()

		return err
	}

	if err := temporary.Close(); err != nil {
		return err
	}

	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}

	removeTemporary = false

	return nil
}
