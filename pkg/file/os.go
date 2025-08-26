package file

import "os"

func Read(path string) (*Source, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return NewSource(path, string(bytes)), nil
}
