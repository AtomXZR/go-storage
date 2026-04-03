//go:build !windows

package local

import "os"

func openReadf(path string) (*os.File, error) {
	file, err := os.Open(path)
	return file, toStorageError(err)
}
