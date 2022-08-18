package core

import (
	"errors"
	"fmt"
)

var errExistsDriver = errors.New("已存在的类型")

type Storages struct {
	storages map[string]Storage
}

func (s *Storages) Get(name string) (Storage, error) {
	if s.storages == nil {
		return nil, errors.New("没有有效的驱动！")
	}

	storage, ok := s.storages[name]
	if !ok {
		return nil, fmt.Errorf("不支持的类型: %s", name)
	}

	return storage, nil
}

func (s *Storages) Register(name string, storage Storage) error {
	if s.storages == nil {
		s.storages = make(map[string]Storage)
	}

	if _, ok := s.storages[name]; ok {
		return errExistsDriver
	}

	s.storages[name] = storage

	return nil
}
