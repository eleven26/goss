package core

import (
	"io"
)

type Store interface {
	Put(key string, r io.Reader) error
	PutFromFile(key string, localPath string) error
	Get(key string) (io.ReadCloser, error)
	SaveToFile(key string, localPath string) error
	Delete(key string) error
	Size(key string) (int64, error)
	Exists(key string) (bool, error)
	Iterator(marker string) FileIterator
}
