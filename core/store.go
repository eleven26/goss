package core

import (
	"io"
)

// Store is an abstraction layer for different cloud storage provider's APIs
type Store interface {
	// Put reads the content of r, uploads it to cloud, and names it as key.
	Put(key string, r io.Reader) error

	// PutFromFile reads the content of localPath, uploads it to cloud, and names it as key.
	PutFromFile(key string, localPath string) error

	// Get gets the content of the object named key in cloud storage.
	Get(key string) (io.ReadCloser, error)

	// Delete deletes the contents of the file pointed to by key.
	Delete(key string) error

	// Size gets the size of the object pointed to by key.
	Size(key string) (int64, error)

	// Exists determines whether the object pointed to by key exists.
	Exists(key string) (bool, error)

	// Iterator gets a FileIterator based on prefix.
	Iterator(prefix string) FileIterator
}
