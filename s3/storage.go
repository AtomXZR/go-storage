package s3

import (
	"context"
	"io"

	"github.com/AtomXZR/go-storage"
	"github.com/minio/minio-go/v7"
)

//
//
//

type S3Storage struct {
	client *minio.Client
	bucket string
}

//
//
//

func (s *S3Storage) Put(ctx context.Context, key string, reader io.Reader, size int64, opts *storage.PutOptions) error {
	key, err := storage.NormalizeKey(key)
	if err != nil {
		return err
	}

	//
	//
	//

	opts = storage.PutOptionsOrDefault(opts)

	_, err = s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		UserMetadata: storage.NormalizeMetadata(opts.Metadata),
		ContentType:  opts.ContentType,
	})

	return s3ErrorToFSError(err)
}

func (s *S3Storage) Get(ctx context.Context, key string, opts *storage.GetOptions) (io.ReadCloser, *storage.Stats, error) {
	key, err := storage.NormalizeKey(key)
	if err != nil {
		return nil, nil, err
	}

	//
	//
	//

	opts = storage.GetOptionsOrDefault(opts)

	getOpts := minio.GetObjectOptions{}

	//

	if opts.Range != nil {
		if err := getOpts.SetRange(opts.Range.Start, opts.Range.End); err != nil {
			return nil, nil, err
		}
	}

	//

	objInfo, err := s.client.StatObject(ctx, s.bucket, key, getOpts) // since StatOpts is just alias...
	if err != nil {
		return nil, nil, s3ErrorToFSError(err)
	}

	reader, err := s.client.GetObject(ctx, s.bucket, key, getOpts)
	if err != nil {
		return nil, nil, s3ErrorToFSError(err)
	}

	stat := objectInfoToStats(objInfo)

	return reader, &stat, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	key, err := storage.NormalizeKey(key)
	if err != nil {
		return err
	}

	//
	//
	//

	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *S3Storage) Stat(ctx context.Context, key string) (*storage.Stats, error) {
	key, err := storage.NormalizeKey(key)
	if err != nil {
		return nil, err
	}

	//
	//
	//

	objInfo, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, s3ErrorToFSError(err)
	}

	stat := objectInfoToStats(objInfo)

	return &stat, nil
}
