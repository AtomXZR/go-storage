package local_test

import (
	"testing"

	"github.com/AtomXZR/go-storage"
	"github.com/AtomXZR/go-storage/local"
	storage_test "github.com/AtomXZR/go-storage/test"
)

var instance storage.Storage

func setup(t *testing.T) (storage.Storage, string) {
	t.Helper()

	// Where `Get` will write file to
	dir := t.TempDir()

	// Where LocalStorage will store the data
	rootDir := t.TempDir()

	if instance == nil {
		inst, err := local.New(rootDir)

		if err != nil {
			t.Fatal(err)
		}

		instance = inst
	}

	return instance, dir
}

func TestStorage(t *testing.T) {
	s, tempdir := setup(t)

	storage_test.DoTest(t, s, tempdir)
}
