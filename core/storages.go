package core

import (
	"fmt"
)

type Storages struct {
	storages map[string]Storage
}

func (s *Storages) Get(name string) Storage {
	if s.storages == nil {
		panic("没有有效的驱动！")
	}

	storage, ok := s.storages[name]
	if !ok {
		panic(fmt.Errorf("不支持的类型: %s", name))
	}

	return storage
}

func (s *Storages) Register(name string, storage Storage) {
	if s.storages == nil {
		s.storages = make(map[string]Storage)
	}

	if _, ok := s.storages[name]; ok {
		panic(fmt.Errorf("已存在的类型：%s", name))
	}

	s.storages[name] = storage
}
