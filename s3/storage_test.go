package s3_test

import (
	"testing"

	"github.com/AtomXZR/go-storage"
	"github.com/AtomXZR/go-storage/s3"
	storage_test "github.com/AtomXZR/go-storage/test"
)

var instance storage.Storage

func setup(t *testing.T) (storage.Storage, string) {
	t.Helper()

	dir := t.TempDir()

	if instance == nil {
		inst, err := s3.New(s3.Config{
			Bucket:    "test-dev",
			Endpoint:  "api.minio.owo-1.home.arpa",
			AccessKey: "mRUzkj9eoGEMmTh5ak3D",
			SecretKey: "9TD9uhx8EtgX3YwD4PofhPUiDRwur7wjzxIwRi1H",

			UseSSL: false,
		})

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
