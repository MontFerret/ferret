package publisher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MontFerret/specs/pkg/api"
)

const moduleID = "montferret/core"

// Publish adds one immutable montferret/core API Reference to an existing Pages tree.
func Publish(pagesRoot string, referenceData []byte) error {
	reference, err := api.Parse(referenceData)
	if err != nil {
		return fmt.Errorf("parse generated API Reference: %w", err)
	}

	if reference.ID != moduleID {
		return fmt.Errorf("generated API Reference id is %q, want %q", reference.ID, moduleID)
	}

	rootInfo, err := os.Lstat(pagesRoot)
	if err != nil {
		return fmt.Errorf("inspect pages root: %w", err)
	}

	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("pages root is not a directory: %s", pagesRoot)
	}

	existing, err := loadExistingState(pagesRoot)
	if err != nil {
		return err
	}

	href := canonicalHref(reference.Version)
	if err := ensureUnpublished(pagesRoot, existing, reference.Version, href); err != nil {
		return err
	}

	index := nextIndex(existing, reference.Version, href)
	if err := api.ValidateIndex(index); err != nil {
		return fmt.Errorf("validate next API Reference index: %w", err)
	}

	indexData, err := encodeIndex(index)
	if err != nil {
		return err
	}

	return writePublication(pagesRoot, reference.Version, referenceData, indexData)
}

func writePublication(pagesRoot, version string, referenceData, indexData []byte) error {
	versionsRoot := filepath.Join(pagesRoot, "versions")
	createdVersionsRoot := false
	if _, err := os.Lstat(versionsRoot); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(versionsRoot, 0o755); err != nil {
			return fmt.Errorf("create versions root: %w", err)
		}

		createdVersionsRoot = true
	} else if err != nil {
		return fmt.Errorf("inspect versions root: %w", err)
	}

	versionDirectory := filepath.Join(versionsRoot, version)
	if err := os.Mkdir(versionDirectory, 0o755); err != nil {
		return fmt.Errorf("create immutable version directory: %w", err)
	}

	rollback := func() {
		_ = os.Remove(filepath.Join(versionDirectory, "api.json"))
		_ = os.Remove(versionDirectory)
		if createdVersionsRoot {
			_ = os.Remove(versionsRoot)
		}
	}

	if err := atomicWrite(filepath.Join(versionDirectory, "api.json"), referenceData); err != nil {
		rollback()

		return fmt.Errorf("write immutable API Reference: %w", err)
	}

	if err := atomicWrite(filepath.Join(pagesRoot, "index.json"), indexData); err != nil {
		rollback()

		return fmt.Errorf("write API Reference index: %w", err)
	}

	return nil
}
