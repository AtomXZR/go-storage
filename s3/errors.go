package s3

import (
	"github.com/AtomXZR/go-storage"
	"github.com/minio/minio-go/v7"
)

func s3ErrorToStorageError(err error) error {
	if err == nil {
		return nil
	}

	resp := minio.ToErrorResponse(err)

	var kind *storage.StorageErrorKind

	switch resp.StatusCode {
	case 400:
		kind = storage.ErrInvalid

	case 403:
		kind = storage.ErrPermission

	case 404:
		kind = storage.ErrNotFound
		if resp.Code == "NoSuchKey" {
			kind = storage.ErrKeyNotExist
		}

	case 405:
		kind = storage.ErrInvalid

	case 409:
		kind = storage.ErrConflict

	case 411:
		kind = storage.ErrInvalid

	case 412:
		kind = storage.ErrInvalid

	case 416:
		kind = storage.ErrInvalid

	case 500:
		kind = storage.ErrStorage

	case 501:
		kind = storage.ErrStorage

	case 503:
		kind = storage.ErrStorage

	case 507:
		kind = storage.ErrStorage

	default:
		kind = storage.ErrUnknown

	}

	return kind.New(err.Error())
}
