package s3_test

import (
	"sync"
	"testing"

	"github.com/AtomXZR/go-storage"
	"github.com/AtomXZR/go-storage/s3"
	storage_test "github.com/AtomXZR/go-storage/test"
	"github.com/joho/godotenv"
)

var instance storage.Storage
var once sync.Once

func setup(t *testing.T) (storage.Storage, string) {
	t.Helper()
	once.Do(func() {
		_ = godotenv.Load()
	})

	dir := t.TempDir()

	if instance == nil {
		inst, err := s3.New(s3.Config{
			Bucket:    storage_test.GetEnvSkip(t, "S3_BUCKET"),
			Endpoint:  storage_test.GetEnvSkip(t, "S3_ENDPOINT"),
			AccessKey: storage_test.GetEnvSkip(t, "S3_ACCESS_KEY"),
			SecretKey: storage_test.GetEnvSkip(t, "S3_SECRET_KEY"),
			UseSSL:    storage_test.GetEnv("S3_USE_SSL") == "true",
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
