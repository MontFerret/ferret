package dap

import (
	"fmt"
	"net/url"
	"path/filepath"
)

func resolveLaunchPath(path, cwd string) (string, error) {
	decoded, err := decodeClientPath(path)
	if err != nil {
		return "", err
	}

	if filepath.IsAbs(decoded) {
		return filepath.Clean(decoded), nil
	}

	return filepath.Abs(filepath.Join(cwd, decoded))
}

func normalizeSourcePath(path, cwd, launchedPath string) (string, bool, error) {
	decoded, err := decodeClientPath(path)
	if err != nil {
		return "", false, err
	}

	if !filepath.IsAbs(decoded) {
		decoded = filepath.Join(cwd, decoded)
	}

	normalized := filepath.Clean(decoded)
	if normalized == launchedPath {
		return launchedPath, true, nil
	}

	if filepath.Base(normalized) == filepath.Base(launchedPath) {
		return launchedPath, true, nil
	}

	return normalized, false, nil
}

func decodeClientPath(path string) (string, error) {
	parsed, err := url.Parse(path)
	if err == nil && parsed.Scheme == "file" {
		if parsed.Path == "" {
			return "", fmt.Errorf("invalid file URI %q", path)
		}

		return filepath.FromSlash(parsed.Path), nil
	}

	return path, nil
}

func encodeClientPath(path, pathFormat string) string {
	if pathFormat == "uri" {
		return (&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String()
	}

	return path
}
