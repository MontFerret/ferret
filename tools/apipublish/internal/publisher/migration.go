package publisher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

// ErrTypeMigrationRequired reports that check mode found legacy string types.
var ErrTypeMigrationRequired = errors.New("core API type migration is required")

type typeMigration struct {
	path        string
	original    []byte
	replacement []byte
}

type migrationReference struct {
	SchemaVersion int                  `json:"schemaVersion"`
	ID            string               `json:"id"`
	Version       string               `json:"version"`
	Namespaces    []migrationNamespace `json:"namespaces"`
}

type migrationNamespace struct {
	Name      string              `json:"name"`
	Functions []migrationFunction `json:"functions"`
}

type migrationFunction struct {
	Name       string               `json:"name"`
	Signatures []migrationSignature `json:"signatures"`
}

type migrationSignature struct {
	Parameters  []migrationParameter `json:"parameters"`
	Variadic    *bool                `json:"variadic,omitempty"`
	Description *string              `json:"description,omitempty"`
	Return      json.RawMessage      `json:"return,omitempty"`
	Throws      *[]api.Throw         `json:"throws,omitempty"`
	Deprecated  *string              `json:"deprecated,omitempty"`
}

type migrationParameter struct {
	Name        string          `json:"name"`
	Type        json.RawMessage `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
}

type migrationReturn struct {
	Type        json.RawMessage `json:"type"`
	Description *string         `json:"description"`
}

// MigrateTypes converts every indexed Core API Reference string type in a
// Pages tree. Check mode validates and reports required replacements without
// writing. Apply mode prepares and validates the complete migration before
// atomically replacing any API Reference.
func MigrateTypes(pagesRoot string, check bool) error {
	index, err := loadMigrationIndex(pagesRoot)
	if err != nil {
		return err
	}

	migrations := make([]typeMigration, 0, len(index.Versions))
	for _, entry := range index.Versions {
		migration, required, err := prepareTypeMigration(pagesRoot, entry)
		if err != nil {
			return err
		}

		if required {
			migrations = append(migrations, migration)
		}
	}

	if len(migrations) == 0 {
		return nil
	}

	if check {
		return fmt.Errorf("%w for %d indexed API Reference(s)", ErrTypeMigrationRequired, len(migrations))
	}

	return applyTypeMigrations(migrations)
}

func loadMigrationIndex(pagesRoot string) (*api.Index, error) {
	rootInfo, err := os.Lstat(pagesRoot)
	if err != nil {
		return nil, fmt.Errorf("inspect pages root: %w", err)
	}

	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("pages root is not a directory: %s", pagesRoot)
	}

	indexPath := filepath.Join(pagesRoot, "index.json")
	indexInfo, err := os.Lstat(indexPath)
	if err != nil {
		return nil, fmt.Errorf("inspect existing index: %w", err)
	}

	if !indexInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("existing index is not a regular file: %s", indexPath)
	}

	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read existing index: %w", err)
	}

	index, err := api.ParseIndex(indexData)
	if err != nil {
		return nil, fmt.Errorf("parse existing API Reference index: %w", err)
	}

	for _, entry := range index.Versions {
		if expected := canonicalHref(entry.Version); entry.Href != expected {
			return nil, fmt.Errorf("existing index version %s has href %q, want authoritative href %q", entry.Version, entry.Href, expected)
		}
	}

	if err := validateVersionTree(pagesRoot, index); err != nil {
		return nil, err
	}

	return index, nil
}

func prepareTypeMigration(pagesRoot string, entry api.IndexVersion) (typeMigration, bool, error) {
	artifactPath := filepath.Join(pagesRoot, filepath.FromSlash(entry.Href[2:]))
	artifactInfo, err := os.Lstat(artifactPath)
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("inspect existing API Reference %s: %w", entry.Version, err)
	}

	if !artifactInfo.Mode().IsRegular() {
		return typeMigration{}, false, fmt.Errorf("existing API Reference is not a regular file: %s", artifactPath)
	}

	artifactData, err := os.ReadFile(artifactPath)
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("read existing API Reference %s: %w", entry.Version, err)
	}

	reference, parseErr := api.Parse(artifactData)
	if parseErr == nil {
		if err := validateMigrationIdentity(reference, entry, artifactPath); err != nil {
			return typeMigration{}, false, err
		}

		if err := validateMigrationCatalog(filepath.Dir(artifactPath), entry.Version, reference); err != nil {
			return typeMigration{}, false, err
		}

		return typeMigration{}, false, nil
	}

	if err := validateSingleJSONDocument(artifactData); err != nil {
		return typeMigration{}, false, fmt.Errorf("parse legacy API Reference %s: %w", entry.Version, err)
	}

	legacy := migrationReference{}
	if err := decodeStrictJSON(artifactData, &legacy); err != nil {
		return typeMigration{}, false, fmt.Errorf("parse legacy API Reference %s: %w", entry.Version, err)
	}

	reference, migrated, err := legacy.structured()
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("convert legacy API Reference %s: %w", entry.Version, err)
	}

	if !migrated {
		return typeMigration{}, false, fmt.Errorf("parse existing API Reference %s: %w", entry.Version, parseErr)
	}

	if err := api.Validate(reference); err != nil {
		return typeMigration{}, false, fmt.Errorf("validate migrated API Reference %s: %w", entry.Version, err)
	}

	replacement, err := encodeReference(reference)
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("encode migrated API Reference %s: %w", entry.Version, err)
	}

	normalized, err := api.Parse(replacement)
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("parse migrated API Reference %s: %w", entry.Version, err)
	}

	if err := validateMigrationIdentity(normalized, entry, artifactPath); err != nil {
		return typeMigration{}, false, err
	}

	if err := validateMigrationCatalog(filepath.Dir(artifactPath), entry.Version, normalized); err != nil {
		return typeMigration{}, false, err
	}

	replacement, err = encodeReference(normalized)
	if err != nil {
		return typeMigration{}, false, fmt.Errorf("encode normalized API Reference %s: %w", entry.Version, err)
	}

	return typeMigration{path: artifactPath, original: artifactData, replacement: replacement}, true, nil
}

func (legacy migrationReference) structured() (*api.Reference, bool, error) {
	migrated := false
	reference := &api.Reference{
		SchemaVersion: legacy.SchemaVersion,
		ID:            legacy.ID,
		Version:       legacy.Version,
		Namespaces:    make([]api.Namespace, len(legacy.Namespaces)),
	}

	for namespaceIndex, namespace := range legacy.Namespaces {
		reference.Namespaces[namespaceIndex] = api.Namespace{
			Name:      namespace.Name,
			Functions: make([]api.Function, len(namespace.Functions)),
		}

		for functionIndex, function := range namespace.Functions {
			targetFunction := &reference.Namespaces[namespaceIndex].Functions[functionIndex]
			targetFunction.Name = function.Name
			targetFunction.Signatures = make([]api.Signature, len(function.Signatures))

			for signatureIndex, signature := range function.Signatures {
				targetSignature := &targetFunction.Signatures[signatureIndex]
				targetSignature.Parameters = make([]api.Parameter, len(signature.Parameters))
				if signature.Variadic != nil {
					if !*signature.Variadic {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d: variadic must be true when present", namespaceIndex, functionIndex, signatureIndex)
					}

					targetSignature.Variadic = true
				}

				if signature.Description != nil {
					if strings.TrimSpace(*signature.Description) == "" {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d: description must not be blank", namespaceIndex, functionIndex, signatureIndex)
					}

					targetSignature.Description = *signature.Description
				}

				if signature.Throws != nil {
					if len(*signature.Throws) == 0 {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d: throws must not be empty", namespaceIndex, functionIndex, signatureIndex)
					}

					targetSignature.Throws = *signature.Throws
				}

				if signature.Deprecated != nil {
					if strings.TrimSpace(*signature.Deprecated) == "" {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d: deprecation must not be blank", namespaceIndex, functionIndex, signatureIndex)
					}

					targetSignature.Deprecated = *signature.Deprecated
				}

				for parameterIndex, parameter := range signature.Parameters {
					parameterType, converted, err := decodeMigrationType(parameter.Type, true)
					if err != nil {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d parameter %d: %w", namespaceIndex, functionIndex, signatureIndex, parameterIndex, err)
					}

					migrated = migrated || converted
					if (parameterType == nil) != (parameter.Description == nil) {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d parameter %d: type and description must appear together", namespaceIndex, functionIndex, signatureIndex, parameterIndex)
					}

					description := ""
					if parameter.Description != nil {
						if strings.TrimSpace(*parameter.Description) == "" {
							return nil, false, fmt.Errorf("namespace %d function %d signature %d parameter %d: description must not be blank", namespaceIndex, functionIndex, signatureIndex, parameterIndex)
						}

						description = *parameter.Description
					}

					targetSignature.Parameters[parameterIndex] = api.Parameter{
						Name:        parameter.Name,
						Type:        parameterType,
						Description: description,
					}
				}

				if len(signature.Return) > 0 {
					returnValue := migrationReturn{}
					if err := decodeStrictJSON(signature.Return, &returnValue); err != nil {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d return: %w", namespaceIndex, functionIndex, signatureIndex, err)
					}

					returnType, converted, err := decodeMigrationType(returnValue.Type, false)
					if err != nil {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d return: %w", namespaceIndex, functionIndex, signatureIndex, err)
					}

					if returnValue.Description == nil || strings.TrimSpace(*returnValue.Description) == "" {
						return nil, false, fmt.Errorf("namespace %d function %d signature %d return: description is required and must not be blank", namespaceIndex, functionIndex, signatureIndex)
					}

					migrated = migrated || converted
					targetSignature.Return = &api.Return{Type: returnType, Description: *returnValue.Description}
				}
			}
		}
	}

	return reference, migrated, nil
}

func decodeMigrationType(data json.RawMessage, optional bool) (*api.Type, bool, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		if optional {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("type is required")
	}

	if bytes.Equal(trimmed, []byte("null")) {
		return nil, false, fmt.Errorf("type must not be null")
	}

	if trimmed[0] == '"' {
		expression := ""
		if err := json.Unmarshal(trimmed, &expression); err != nil {
			return nil, false, fmt.Errorf("decode legacy type expression: %w", err)
		}

		parsed, err := api.ParseType(expression)
		if err != nil {
			return nil, false, fmt.Errorf("parse legacy type expression %q: %w", expression, err)
		}

		return &parsed, true, nil
	}

	structured := api.Type{}
	if err := decodeStrictJSON(trimmed, &structured); err != nil {
		return nil, false, fmt.Errorf("decode structured type: %w", err)
	}

	return &structured, false, nil
}

func validateMigrationIdentity(reference *api.Reference, entry api.IndexVersion, path string) error {
	if reference.ID != moduleID || reference.Version != entry.Version {
		return fmt.Errorf(
			"existing API Reference %s identifies %s@%s, want %s@%s",
			path,
			reference.ID,
			reference.Version,
			moduleID,
			entry.Version,
		)
	}

	return nil
}

func validateMigrationCatalog(directory, version string, reference *api.Reference) error {
	catalogPath := filepath.Join(directory, "catalog.json")
	catalogInfo, err := os.Lstat(catalogPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("inspect existing API Catalog %s: %w", version, err)
	}

	if !catalogInfo.Mode().IsRegular() {
		return fmt.Errorf("existing API Catalog is not a regular file: %s", catalogPath)
	}

	catalogData, err := os.ReadFile(catalogPath)
	if err != nil {
		return fmt.Errorf("read existing API Catalog %s: %w", version, err)
	}

	catalog, err := apicatalog.Parse(catalogData)
	if err != nil {
		return fmt.Errorf("parse existing API Catalog %s: %w", version, err)
	}

	if err := validatePair(reference, catalog); err != nil {
		return fmt.Errorf("validate existing API Reference and Catalog %s: %w", version, err)
	}

	return nil
}

func encodeReference(reference *api.Reference) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(reference); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func applyTypeMigrations(migrations []typeMigration) error {
	written := make([]typeMigration, 0, len(migrations))
	for _, migration := range migrations {
		if err := atomicWrite(migration.path, migration.replacement); err != nil {
			rollbackErrors := make([]error, 0, len(written))
			for index := len(written) - 1; index >= 0; index-- {
				previous := written[index]
				if rollbackErr := atomicWrite(previous.path, previous.original); rollbackErr != nil {
					rollbackErrors = append(rollbackErrors, fmt.Errorf("restore %s: %w", previous.path, rollbackErr))
				}
			}

			return errors.Join(fmt.Errorf("write migrated API Reference %s: %w", migration.path, err), errors.Join(rollbackErrors...))
		}

		written = append(written, migration)
	}

	return nil
}

func validateSingleJSONDocument(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := consumeJSONValue(decoder); err != nil {
		return err
	}

	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value")
		}

		return err
	}

	return nil
}

func consumeJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	if token == nil {
		return fmt.Errorf("null values are not allowed")
	}

	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delimiter {
	case '{':
		keys := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}

			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object key is not a string")
			}

			if _, exists := keys[key]; exists {
				return fmt.Errorf("duplicate object key %q", key)
			}

			keys[key] = struct{}{}
			if err := consumeJSONValue(decoder); err != nil {
				return err
			}
		}

		closing, err := decoder.Token()
		if err != nil {
			return err
		}

		if closing != json.Delim('}') {
			return fmt.Errorf("object has invalid closing delimiter")
		}
	case '[':
		for decoder.More() {
			if err := consumeJSONValue(decoder); err != nil {
				return err
			}
		}

		closing, err := decoder.Token()
		if err != nil {
			return err
		}

		if closing != json.Delim(']') {
			return fmt.Errorf("array has invalid closing delimiter")
		}
	default:
		return fmt.Errorf("unexpected closing delimiter %q", delimiter)
	}

	return nil
}

func decodeStrictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value")
		}

		return err
	}

	return nil
}
