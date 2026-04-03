package local

import (
	"errors"
	"os"

	"github.com/AtomXZR/go-storage"
)

func toStorageError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, os.ErrPermission):
		return storage.ErrPermission.New(err.Error())
	case errors.Is(err, os.ErrInvalid):
		return storage.ErrInvalid.New(err.Error())
	default:
		return storage.ErrStorage.New(err.Error())
	}
}
