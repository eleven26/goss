package huawei

import (
	"io"
	"io/ioutil"
	"strconv"

	"github.com/eleven26/goss/core"
)

type Storage struct {
	store Store
}

func (s *Storage) Put(key string, r io.Reader) error {
	return s.store.Put(key, r)
}

func (s *Storage) PutFromFile(key string, localPath string) error {
	return s.store.PutFromFile(key, localPath)
}

func (s *Storage) Get(key string) (io.ReadCloser, error) {
	return s.store.Get(key)
}

func (s *Storage) GetString(key string) (string, error) {
	bs, err := s.GetBytes(key)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (s *Storage) GetBytes(key string) ([]byte, error) {
	rc, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	defer func(rc io.ReadCloser) {
		err = rc.Close()
	}(rc)

	return ioutil.ReadAll(rc)
}

func (s *Storage) GetToFile(key string, localPath string) error {
	return s.store.SaveToFile(key, localPath)
}

func (s *Storage) Delete(key string) error {
	return s.store.Delete(key)
}

func (s *Storage) Exists(key string) (bool, error) {
	return s.store.Exists(key)
}

func (s *Storage) Files(dir string) ([]core.File, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Size(key string) (int64, error) {
	header, err := s.store.Meta(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(header.Get("Content-Length"), 10, 64)
}

func (s *Storage) Storage() interface{} {
	return s
}
