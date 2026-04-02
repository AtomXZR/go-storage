package local

import "path/filepath"

func New(rootDir string) (*LocalStorage, error) {
	path := filepath.Clean(rootDir)
	abs, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	return &LocalStorage{
		rootDir: abs,
	}, nil
}
