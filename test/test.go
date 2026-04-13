package test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/AtomXZR/go-storage"
)

//
//
//

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvSkip(t *testing.T, key string) string {
	t.Helper()
	v := GetEnv(key)
	if v == "" {
		t.Skipf("skipping: %s not set", key)
	}
	return v
}

//
//
//

func dirname() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func getPath(fn string) string {
	return filepath.Join(dirname(), fn)
}

func OpenFile(fn string) (*os.File, error) {
	return os.Open(getPath(fn))
}

//
//
//

func GetTextFile() (*os.File, error) {
	return OpenFile("text.txt")
}

func GetImageFile() (*os.File, error) {
	return OpenFile("image.jpg")
}

//
//
//

func DoTest(t *testing.T, store storage.Storage, tempDir string) {
	t.Helper()

	//

	t.Run("Put", func(t *testing.T) {
		t.Run("invalid key", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "", f, fi.Size(), nil)

			if !errors.Is(err, storage.ErrInvalidKey) {
				t.Fatalf("expected ErrInvalidKey, got %v", err)
			}
		})

		t.Run("text file", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			//

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "text/hello.txt", f, fi.Size(), nil)
			if err != nil {
				t.Fatalf("Put: %v", err)
			}
		})

		t.Run("image file", func(t *testing.T) {
			f, err := GetImageFile()
			if err != nil {
				t.Fatalf("GetImageFile: %v", err)
			}
			defer f.Close()

			//

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "images/photo.png", f, fi.Size(), nil)
			if err != nil {
				t.Fatalf("Put: %v", err)
			}
		})

		t.Run("with content type", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			//

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "text/typed.txt", f, fi.Size(), &storage.PutOptions{
				ContentType: "text/plain",
			})
			if err != nil {
				t.Fatalf("Put with content type: %v", err)
			}
		})

		t.Run("with metadata", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			//

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "text/with-meta.txt", f, fi.Size(), &storage.PutOptions{
				Metadata: storage.Metadata{
					"author": "test",
					"env":    "ci",
				},
			})
			if err != nil {
				t.Fatalf("Put with metadata: %v", err)
			}
		})

		t.Run("overwrite existing key", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			//

			fi, _ := f.Stat()

			err = store.Put(t.Context(), "text/hello.txt", f, fi.Size(), nil)
			if err != nil {
				t.Fatalf("Put overwrite: %v", err)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("invalid key", func(t *testing.T) {
			rc, stats, err := store.Get(t.Context(), "", nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, storage.ErrInvalidKey) {
				t.Fatalf("expected ErrInvalidKey, got %v", err)
			}

			if rc != nil {
				t.Error("expected nil reader")
			}
			if stats != nil {
				t.Error("expected nil stats")
			}
		})

		t.Run("existing key", func(t *testing.T) {
			rc, stats, err := store.Get(t.Context(), "text/hello.txt", nil)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			defer rc.Close()

			//

			if stats == nil {
				t.Fatal("expected stats, got nil")
			}
			if stats.Size == 0 {
				t.Error("expected non-zero size")
			}
		})

		t.Run("content matches original", func(t *testing.T) {
			f, err := GetTextFile()
			if err != nil {
				t.Fatalf("GetTextFile: %v", err)
			}
			defer f.Close()

			//

			original, err := io.ReadAll(f)
			if err != nil {
				t.Fatalf("read original: %v", err)
			}

			rc, _, err := store.Get(t.Context(), "text/hello.txt", nil)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			defer rc.Close()

			//

			got, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("read response: %v", err)
			}

			if !bytes.Equal(original, got) {
				t.Errorf("content mismatch: got %d bytes, want %d bytes", len(got), len(original))
			}
		})

		t.Run("content type is preserved", func(t *testing.T) {
			_, stats, err := store.Get(t.Context(), "text/typed.txt", nil)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}

			if stats.ContentType != "text/plain" {
				t.Errorf("ContentType = %q, want %q", stats.ContentType, "text/plain")
			}
		})

		t.Run("metadata is preserved", func(t *testing.T) {
			_, stats, err := store.Get(t.Context(), "text/with-meta.txt", nil)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}

			t.Logf("META: %v", stats)

			if stats.Metadata["Author"] != "test" {
				t.Errorf("Metadata[Author] = %q, want %q", stats.Metadata["Author"], "test")
			}
			if stats.Metadata["Env"] != "ci" {
				t.Errorf("Metadata[Env] = %q, want %q", stats.Metadata["Env"], "ci")
			}
		})

		t.Run("byte range", func(t *testing.T) {
			rc, stats, err := store.Get(t.Context(), "text/hello.txt", &storage.GetOptions{
				Range: &storage.Range{
					Start: 0,
					End:   3,
				},
			})
			if err != nil {
				t.Fatalf("Get with range: %v", err)
			}
			defer rc.Close()

			got, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("read range response: %v", err)
			}

			if int64(len(got)) != stats.Size {
				t.Errorf("range read: got %d bytes, stats.Size = %d", len(got), stats.Size)
			}
			if len(got) != 4 { // 0-3 inclusive = 4 bytes
				t.Errorf("range read: got %d bytes, want 4", len(got))
			}
		})

		t.Run("non-existent key", func(t *testing.T) {
			_, _, err := store.Get(t.Context(), "does/not/exist.txt", nil)
			if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("expected os.ErrNotExist, got %v", err)
			}
		})
	})

	t.Run("Stat", func(t *testing.T) {
		t.Run("existing key", func(t *testing.T) {
			stats, err := store.Stat(t.Context(), "text/hello.txt")
			if err != nil {
				t.Fatalf("Stat: %v", err)
			}

			if stats.Size == 0 {
				t.Error("expected non-zero size")
			}
			if stats.LastModified.IsZero() {
				t.Error("expected non-zero LastModified")
			}
			if stats.ETag == "" {
				t.Error("expected non-empty ETag")
			}
		})

		t.Run("matches Get stats", func(t *testing.T) {
			statStats, err := store.Stat(t.Context(), "text/hello.txt")
			if err != nil {
				t.Fatalf("Stat: %v", err)
			}

			_, getStats, err := store.Get(t.Context(), "text/hello.txt", nil)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}

			if statStats.Size != getStats.Size {
				t.Errorf("Size mismatch: Stat=%d Get=%d", statStats.Size, getStats.Size)
			}
			if statStats.ETag != getStats.ETag {
				t.Errorf("ETag mismatch: Stat=%q Get=%q", statStats.ETag, getStats.ETag)
			}
		})

		t.Run("non-existent key", func(t *testing.T) {
			_, err := store.Stat(t.Context(), "does/not/exist.txt")
			if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("expected os.ErrNotExist, got %v", err)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("existing key", func(t *testing.T) {
			err := store.Delete(t.Context(), "text/hello.txt")
			if err != nil {
				t.Fatalf("Delete: %v", err)
			}

			_, err = store.Stat(t.Context(), "text/hello.txt")
			if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("expected ErrNotExist after delete, got %v", err)
			}
		})

		t.Run("non-existent key", func(t *testing.T) {
			// Should not error - idempotent delete
			err := store.Delete(t.Context(), "does/not/exist.txt")
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Errorf("unexpected error deleting non-existent key: %v", err)
			}
		})
	})
}
