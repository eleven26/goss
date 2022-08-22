package core

import (
	"io"
	"net/http"
)

type Store interface {
	Put(key string, r io.Reader) error
	PutFromFile(key string, localPath string) error
	Get(key string) (io.ReadCloser, error)
	SaveToFile(key string, localPath string) error
	Delete(key string) error
	Meta(key string) (http.Header, error)
	Exists(key string) (bool, error)
	Iterator(marker string) FileIterator
}
