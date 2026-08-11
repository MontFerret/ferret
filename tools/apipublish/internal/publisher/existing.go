package publisher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MontFerret/specs/pkg/api"
)

func loadExistingState(pagesRoot string) (*api.Index, error) {
	indexPath := filepath.Join(pagesRoot, "index.json")
	info, err := os.Lstat(indexPath)
	if errors.Is(err, os.ErrNotExist) {
		if err := validateVersionTree(pagesRoot, nil); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("inspect existing index: %w", err)
	}

	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("existing index is not a regular file: %s", indexPath)
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read existing index: %w", err)
	}

	index, err := api.ParseIndex(data)
	if err != nil {
		return nil, fmt.Errorf("parse existing API Reference index: %w", err)
	}

	for _, entry := range index.Versions {
		expectedHref := canonicalHref(entry.Version)
		if entry.Href != expectedHref {
			return nil, fmt.Errorf("existing index version %s has href %q, want authoritative href %q", entry.Version, entry.Href, expectedHref)
		}

		artifactPath := filepath.Join(pagesRoot, filepath.FromSlash(entry.Href[2:]))
		artifactInfo, err := os.Lstat(artifactPath)
		if err != nil {
			return nil, fmt.Errorf("inspect existing API Reference %s: %w", entry.Version, err)
		}

		if !artifactInfo.Mode().IsRegular() {
			return nil, fmt.Errorf("existing API Reference is not a regular file: %s", artifactPath)
		}

		artifactData, err := os.ReadFile(artifactPath)
		if err != nil {
			return nil, fmt.Errorf("read existing API Reference %s: %w", entry.Version, err)
		}

		reference, err := api.Parse(artifactData)
		if err != nil {
			return nil, fmt.Errorf("parse existing API Reference %s: %w", entry.Version, err)
		}

		if reference.ID != moduleID || reference.Version != entry.Version {
			return nil, fmt.Errorf(
				"existing API Reference %s identifies %s@%s, want %s@%s",
				artifactPath,
				reference.ID,
				reference.Version,
				moduleID,
				entry.Version,
			)
		}
	}

	if err := validateVersionTree(pagesRoot, index); err != nil {
		return nil, err
	}

	return index, nil
}

func validateVersionTree(pagesRoot string, index *api.Index) error {
	versionsRoot := filepath.Join(pagesRoot, "versions")
	info, err := os.Lstat(versionsRoot)
	if errors.Is(err, os.ErrNotExist) {
		if index != nil {
			return fmt.Errorf("existing index references artifacts but versions directory does not exist")
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("inspect versions directory: %w", err)
	}

	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("versions path is not a real directory: %s", versionsRoot)
	}

	entries, err := os.ReadDir(versionsRoot)
	if err != nil {
		return fmt.Errorf("read versions directory: %w", err)
	}

	expected := make(map[string]struct{})
	if index != nil {
		for _, version := range index.Versions {
			expected[version.Version] = struct{}{}
		}
	}

	for _, entry := range entries {
		if _, exists := expected[entry.Name()]; !exists {
			return fmt.Errorf("unindexed API Reference version exists: %s", filepath.Join(versionsRoot, entry.Name()))
		}

		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect API Reference version %s: %w", entry.Name(), err)
		}

		if !entryInfo.IsDir() || entryInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("API Reference version is not a real directory: %s", filepath.Join(versionsRoot, entry.Name()))
		}

		versionEntries, err := os.ReadDir(filepath.Join(versionsRoot, entry.Name()))
		if err != nil {
			return fmt.Errorf("read API Reference version %s: %w", entry.Name(), err)
		}

		if len(versionEntries) != 1 || versionEntries[0].Name() != "api.json" {
			return fmt.Errorf("API Reference version %s must contain only api.json", entry.Name())
		}
	}

	if len(entries) != len(expected) {
		return fmt.Errorf("versions directory contains %d entries but index contains %d", len(entries), len(expected))
	}

	return nil
}

func ensureUnpublished(pagesRoot string, existing *api.Index, version, href string) error {
	versionDirectory := filepath.Join(pagesRoot, "versions", version)
	if _, err := os.Lstat(versionDirectory); err == nil {
		return fmt.Errorf("immutable API Reference version directory already exists: %s", versionDirectory)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect target version directory: %w", err)
	}

	if existing == nil {
		return nil
	}

	for _, entry := range existing.Versions {
		if entry.Version == version {
			return fmt.Errorf("immutable API Reference version is already indexed: %s", version)
		}

		if entry.Href == href {
			return fmt.Errorf("immutable API Reference href is already indexed: %s", href)
		}
	}

	return nil
}
