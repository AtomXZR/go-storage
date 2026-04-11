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

	//

	if opts.Range != nil {
		start := opts.Range.Start
		end := opts.Range.End

		//

		if end == -1 {
			end = meta.Size - 1
		}

		//

		if start < 0 || end < start || end >= meta.Size {
			return nil, nil, storage.ErrInvalid.New("invalid range")
		}

		length := end - start + 1

		file, err := openReadf(dataFilePath)
		if err != nil {
			return nil, nil, err
		}

		section := io.NewSectionReader(file, start, length)

		//

		outMeta := *meta
		outMeta.Size = length

		return struct {
			io.Reader
			io.Closer
		}{
			Reader: section,
			Closer: file,
		}, &outMeta, nil
	}

	//

	file, err := openReadf(dataFilePath)
	if err != nil {
		return nil, nil, err
	}

	return file, meta, nil
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
