package local

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/AtomXZR/go-storage"
	"github.com/cespare/xxhash/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//
//
//

func keyToPath(key string) string {
	if key == "" {
		return ""
	}

	hash := xxhash.Sum64String(key)

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], hash)

	//

	h := hex.EncodeToString(buf[:])

	seg1 := h[0:2]
	seg2 := h[2:4]
	rest := h[4:]

	return filepath.Join(seg1, seg2, rest)
}

//

func mkDirAll(path string) error {
	return toStorageError(os.MkdirAll(path, 0755))
}

func isDirExist(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	return f.IsDir()
}

// func isFileExist(path string) bool {
// 	f, err := os.Stat(path)
// 	if err != nil {
// 		return false
// 	}

// 	return !f.IsDir()
// }

func createf(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	return file, toStorageError(err)
}

func readAll(path string) ([]byte, error) {
	file, err := openReadf(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, err := io.ReadAll(file)
	if err != nil {
		return nil, toStorageError(err)
	}

	return result, nil
}

//
//
//

func writeDataFile(path string, r io.Reader, size int64) (hash string, err error) {
	file, err := createf(path)
	if err != nil {
		return "", toStorageError(err)
	}
	defer file.Close()

	hasher := sha256.New()

	multiWriter := io.MultiWriter(file, hasher)
	if _, err := io.CopyN(multiWriter, r, size); err != nil {
		return "", toStorageError(err)
	}

	//

	if err := file.Sync(); err != nil {
		return "", toStorageError(err)
	}

	//

	hash = hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}

func statsToJson(stats storage.Stats) (result string, err error) {
	result, _ = sjson.Set("", "ETag", stats.ETag)
	result, _ = sjson.Set(result, "Content-Type", stats.ContentType)
	result, _ = sjson.Set(result, "Size", stats.Size)
	result, _ = sjson.Set(result, "Last-Modified", stats.LastModified)
	result, _ = sjson.Set(result, "Metadata", stats.Metadata)

	return result, nil
}

func jsonToStats(json []byte) (result *storage.Stats, err error) {
	meta := gjson.ParseBytes(json)

	//

	metadata := make(storage.Metadata)
	meta.Get("Metadata").ForEach(func(key, value gjson.Result) bool {
		metadata[storage.NormalizeMetadataKey(key.String())] = value.String()
		return true
	})

	//

	result = &storage.Stats{
		ETag:         meta.Get("ETag").String(),
		ContentType:  meta.Get("Content-Type").String(),
		Size:         meta.Get("Size").Int(),
		LastModified: meta.Get("Last-Modified").Time(),

		Metadata: metadata,
	}

	return result, nil
}

func writeMetadataFile(path string, stats storage.Stats) (err error) {
	file, err := createf(path)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := statsToJson(stats)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(json); err != nil {
		return toStorageError(err)
	}

	return toStorageError(file.Sync())
}

func readMetadataFile(path string) (*storage.Stats, error) {
	buf, err := readAll(path)
	if err != nil {
		return nil, err
	}

	return jsonToStats(buf)
}
