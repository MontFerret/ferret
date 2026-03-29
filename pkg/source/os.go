package source

import "os"

func Read(path string) (*Source, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return New(path, string(bytes)), nil
}
