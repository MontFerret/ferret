package publisher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

const moduleID = "montferret/core"

// Publish adds one immutable montferret/core API Reference and Catalog to an existing Pages tree.
func Publish(pagesRoot string, referenceData, catalogData []byte) error {
	reference, err := api.Parse(referenceData)
	if err != nil {
		return fmt.Errorf("parse generated API Reference: %w", err)
	}

	if reference.ID != moduleID {
		return fmt.Errorf("generated API Reference id is %q, want %q", reference.ID, moduleID)
	}

	catalog, err := apicatalog.Parse(catalogData)
	if err != nil {
		return fmt.Errorf("parse generated API Catalog: %w", err)
	}

	if err := validatePair(reference, catalog); err != nil {
		return fmt.Errorf("validate generated API Reference and Catalog: %w", err)
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

	return writePublication(pagesRoot, reference.Version, referenceData, catalogData, indexData)
}

func writePublication(pagesRoot, version string, referenceData, catalogData, indexData []byte) error {
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

	stagingDirectory, err := os.MkdirTemp(versionsRoot, ".api-publication-"+version+"-*")
	if err != nil {
		if createdVersionsRoot {
			_ = os.Remove(versionsRoot)
		}

		return fmt.Errorf("create publication staging directory: %w", err)
	}
	removeStaging := true
	publicationComplete := false
	defer func() {
		if removeStaging {
			removePublicationDirectory(stagingDirectory)
		}
		if createdVersionsRoot && !publicationComplete {
			_ = os.Remove(versionsRoot)
		}
	}()

	if err := os.Chmod(stagingDirectory, 0o755); err != nil {
		return fmt.Errorf("set publication staging permissions: %w", err)
	}

	if err := atomicWrite(filepath.Join(stagingDirectory, "api.json"), referenceData); err != nil {
		return fmt.Errorf("write staged API Reference: %w", err)
	}

	if err := atomicWrite(filepath.Join(stagingDirectory, "catalog.json"), catalogData); err != nil {
		return fmt.Errorf("write staged API Catalog: %w", err)
	}

	versionDirectory := filepath.Join(versionsRoot, version)
	if err := os.Rename(stagingDirectory, versionDirectory); err != nil {
		return fmt.Errorf("publish immutable API artifact directory: %w", err)
	}
	removeStaging = false

	rollback := func() {
		removePublicationDirectory(versionDirectory)
		if createdVersionsRoot {
			_ = os.Remove(versionsRoot)
		}
	}

	if err := atomicWrite(filepath.Join(pagesRoot, "index.json"), indexData); err != nil {
		rollback()

		return fmt.Errorf("write API Reference index: %w", err)
	}
	publicationComplete = true

	return nil
}

func removePublicationDirectory(directory string) {
	_ = os.Remove(filepath.Join(directory, "catalog.json"))
	_ = os.Remove(filepath.Join(directory, "api.json"))
	_ = os.Remove(directory)
}
