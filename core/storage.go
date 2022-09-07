package core

import (
	"io"
	"os"
)

// Storage defines a unified interface for reading and writing cloud storage objects.
type Storage interface {
	// Put saves the content read from r to the key of oss.
	Put(key string, r io.Reader) error

	// PutFromFile saves the file pointed to by the `localPath` to the oss key.
	PutFromFile(key string, localPath string) error

	// Get gets the file pointed to by key.
	Get(key string) (io.ReadCloser, error)

	// GetString gets the file pointed to by key and returns a string.
	GetString(key string) (string, error)

	// GetBytes gets the file pointed to by key and returns a byte array.
	GetBytes(key string) ([]byte, error)

	// GetToFile saves the file pointed to by key to the localPath.
	GetToFile(key string, localPath string) error

	// Delete the file pointed to by key.
	Delete(key string) error

	// Exists determines whether the file exists.
	Exists(key string) (bool, error)

	// Files list all files in the specified directory.
	Files(dir string) ([]File, error)

	// Size fet the file size.
	Size(key string) (int64, error)

	// Store is an instance for calling APIs of different cloud storage service providers.
	Store() interface{}
}

type storage struct {
	store Store
}

// NewStorage create new Storage instance.
func NewStorage(store Store) Storage {
	return &storage{
		store: store,
	}
}

// Put saves the content read from r to the key of oss.
func (s *storage) Put(key string, r io.Reader) error {
	return s.store.Put(key, r)
}

// PutFromFile saves the file pointed to by the `localPath` to the oss key.
func (s *storage) PutFromFile(key string, localPath string) error {
	return s.store.PutFromFile(key, localPath)
}

// Get gets the file pointed to by key.
func (s *storage) Get(key string) (reader io.ReadCloser, err error) {
	return s.store.Get(key)
}

// GetString gets the file pointed to by key and returns a string.
func (s *storage) GetString(key string) (string, error) {
	bs, err := s.GetBytes(key)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// GetBytes gets the file pointed to by key and returns a byte array.
func (s *storage) GetBytes(key string) (bytes []byte, err error) {
	rc, err := s.store.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	return io.ReadAll(rc)
}

// GetToFile saves the file pointed to by key to the localPath.
func (s *storage) GetToFile(key string, localPath string) (err error) {
	rc, err := s.store.Get(key)
	if err != nil {
		return err
	}

	defer func(rc io.ReadCloser) {
		err = rc.Close()
	}(rc)

	f, _ := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	_, err = io.Copy(f, rc)

	return err
}

// Size fet the file size.
func (s *storage) Size(key string) (int64, error) {
	return s.store.Size(key)
}

// Delete the file pointed to by key.
func (s *storage) Delete(key string) error {
	return s.store.Delete(key)
}

// Exists determines whether the file exists.
func (s *storage) Exists(key string) (bool, error) {
	return s.store.Exists(key)
}

// Files list all files in the specified directory.
func (s *storage) Files(dir string) ([]File, error) {
	return s.store.Iterator(dir).All()
}

func (s *storage) Store() interface{} {
	return s.store
}
