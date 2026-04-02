package storage

import (
	"context"
	"io"
	"time"
)

// Key: /[a-z0-9\-_]+/
type Metadata map[string]string

type Stats struct {
	ETag         string
	ContentType  string
	Size         int64
	LastModified time.Time

	Metadata Metadata
}

// zero-indexed & inclusive; https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Range
type Range struct {
	Start int64
	End   int64
}

type GetOptions struct {
	Range *Range
}

type PutOptions struct {
	ContentType string
	Metadata    Metadata
}

type Storage interface {
	Put(ctx context.Context, key string, reader io.Reader, size int64, opts *PutOptions) error

	Get(ctx context.Context, key string, opts *GetOptions) (io.ReadCloser, *Stats, error)

	Delete(ctx context.Context, key string) error

	Stat(ctx context.Context, key string) (*Stats, error)
}
