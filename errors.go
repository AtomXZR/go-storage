package storage

import "errors"

//
//
//

type StorageErrorKind struct {
	name string
}

func (k *StorageErrorKind) Error() string {
	return k.name
}

func (k *StorageErrorKind) New(msg string) StorageError {
	return StorageError{
		kind:    k,
		message: msg,
	}
}

func NewStorageErrorKind(name string) *StorageErrorKind {
	return &StorageErrorKind{name: name}
}

//
//
//

type StorageError struct {
	kind    *StorageErrorKind
	message string
}

func (e StorageError) Error() string {
	if e.message == "" {
		return e.kind.name
	}
	return e.kind.name + ": " + e.message
}

func (e StorageError) Unwrap() error {
	return e.kind
}

//
//
//

func AsStorageError(err error) (StorageError, bool) {
	var se StorageError
	ok := errors.As(err, &se)
	return se, ok
}

//
//
//

var ErrUnknown = NewStorageErrorKind("unknown error")

var ErrStorage = NewStorageErrorKind("storage error") // HTTP 500
var ErrInvalid = NewStorageErrorKind("invalid argument")
var ErrConflict = NewStorageErrorKind("conflict state")
var ErrNotFound = NewStorageErrorKind("not found")

var ErrInvalidKey = NewStorageErrorKind("invalid key")
var ErrKeyNotExist = NewStorageErrorKind("key does not exist")

var ErrPermission = NewStorageErrorKind("permission denied")
