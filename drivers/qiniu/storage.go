package qiniu

import (
	"io"
	"io/ioutil"

	"github.com/eleven26/goss/core"
)

type Storage struct {
	store Store
}

func (s *Storage) Put(key string, r io.Reader) error {
	_, err := s.store.put(key, r)

	return err
}

func (s *Storage) PutFromFile(key string, localPath string) error {
	_, err := s.store.putFromFile(key, localPath)

	return err
}

func (s *Storage) Get(key string) (io.ReadCloser, error) {
	url := s.store.getDownloadUrl(key)

	return s.store.getUrlContent(url)
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

func (s *Storage) Delete(key string) error {
	return s.store.Delete(key)
}

func (s *Storage) Save(key string, localPath string) error {
	return s.store.SaveToFile(key, localPath)
}

func (s *Storage) Exists(key string) (bool, error) {
	_, err := s.store.Stat(key)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *Storage) Files(dir string) ([]core.File, error) {
	return s.store.Iterator(dir).All()
}

func (s *Storage) Size(key string) (int64, error) {
	fileInfo, err := s.store.Stat(key)
	if err != nil {
		return 0, err
	}

	return fileInfo.Fsize, nil
}

func (s *Storage) Storage() interface{} {
	return s
}
