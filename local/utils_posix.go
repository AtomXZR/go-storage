//go:build !windows

package local

import "os"

func openReadf(path string) (*os.File, error) {
	return os.Open(path)
}
