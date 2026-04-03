package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AtomXZR/go-storage"
)

type LocalStorage struct {
	rootDir string
}

//
//
//

func (s *LocalStorage) getPath(key string) (baseDir string, dataFilePath string, metaFilePath string, err error) {
	key, err = storage.NormalizeKey(key)
	if err != nil {
		return "", "", "", err // InvalidKey error
	}

	//

	segment := keyToPath(key)

	baseDir = filepath.Join(s.rootDir, segment)
	dataFilePath = filepath.Join(baseDir, "data.bin")
	metaFilePath = filepath.Join(baseDir, "metadata.json")

	return baseDir, dataFilePath, metaFilePath, nil
}

//
//
//

func (s *LocalStorage) Put(ctx context.Context, key string, reader io.Reader, size int64, opts *storage.PutOptions) error {
	opts = storage.PutOptionsOrDefault(opts)
	baseDir, dataFilePath, metaFilePath, err := s.getPath(key)
	if err != nil {
		return err
	}

	if err := mkDirAll(baseDir); err != nil {
		return err
	}

	//

	tempDataFilePath := dataFilePath + ".tmp"
	tempMetaFilePath := metaFilePath + ".tmp"

	//

	hash, err := writeDataFile(tempDataFilePath, reader, size)
	if err != nil {
		return err
	}

	err = writeMetadataFile(tempMetaFilePath, storage.Stats{
		ETag:         hash,
		ContentType:  opts.ContentType,
		Size:         size,
		LastModified: time.Now(),

		Metadata: opts.Metadata,
	})
	if err != nil {
		return err
	}

	//

	if err := os.Rename(tempDataFilePath, dataFilePath); err != nil {
		return err
	}

	if err := os.Rename(tempMetaFilePath, metaFilePath); err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Get(ctx context.Context, key string, opts *storage.GetOptions) (io.ReadCloser, *storage.Stats, error) {
	opts = storage.GetOptionsOrDefault(opts)
	baseDir, dataFilePath, metaFilePath, err := s.getPath(key)
	if err != nil {
		return nil, nil, err
	}

	if !isDirExist(baseDir) {
		return nil, nil, os.ErrNotExist
	}

	//

	meta, err := readMetadataFile(metaFilePath)
	if err != nil {
		return nil, nil, err
	}

	file, err := openReadf(dataFilePath)
	if err != nil {
		return nil, nil, err
	}

	var rc io.ReadCloser = file

	//

	if opts.Range != nil {
		length := opts.Range.End - opts.Range.Start + 1

		section := io.NewSectionReader(file, opts.Range.Start, length)
		rc = struct {
			io.Reader
			io.Closer
		}{section, file}

		meta.Size = length
	}

	//

	return rc, meta, nil
}

func (s *LocalStorage) Stat(ctx context.Context, key string) (*storage.Stats, error) {
	baseDir, _, metaFilePath, err := s.getPath(key)
	if err != nil {
		return nil, err
	}

	if !isDirExist(baseDir) {
		return nil, os.ErrNotExist
	}

	//

	return readMetadataFile(metaFilePath)
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	baseDir, _, _, err := s.getPath(key)
	if err != nil {
		return err
	}

	if !isDirExist(baseDir) {
		return nil
	}

	//

	return os.RemoveAll(baseDir)
}
