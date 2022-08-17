package qiniu

import (
	"github.com/eleven26/goss/core"
)

type Storage struct {
	store Store
}

func (s *Storage) Put(key string, localPath string) error {
	_, err := s.store.put(key, localPath)

	return err
}

func (s *Storage) Get(key string) (string, error) {
	url := s.store.getDownloadUrl(key)

	bs, err := s.store.getUrlContent(url)
	if err != nil {
		return "", err
	}

	return string(bs), nil
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
