package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StorageStub struct{}

func (s StorageStub) Put(key string, localPath string) error {
	panic("implement me")
}

func (s StorageStub) Get(key string) (string, error) {
	panic("implement me")
}

func (s StorageStub) Delete(key string) error {
	panic("implement me")
}

func (s StorageStub) Save(key string, localPath string) error {
	panic("implement me")
}

func (s StorageStub) Exists(key string) (bool, error) {
	panic("implement me")
}

func (s StorageStub) Files(dir string) ([]File, error) {
	panic("implement me")
}

func (s StorageStub) Storage() interface{} {
	panic("implement me")
}

func (s StorageStub) Size(key string) (int64, error) {
	panic("implement me")
}

func TestRegister(t *testing.T) {
	storages := Storages{}

	assert.Nil(t, storages.storages)
	stub := StorageStub{}
	storages.Register("test", stub)
	assert.NotNil(t, storages.storages)

	assert.Len(t, storages.storages, 1)

	s, ok := storages.storages["test"]
	assert.True(t, ok)
	assert.Equal(t, s, stub)

	assert.Panics(t, func() {
		storages.Register("test", stub)
	})

	storages.Register("test1", StorageStub{})
	assert.Len(t, storages.storages, 2)
}

func TestGet(t *testing.T) {
	storages := Storages{}

	assert.Panics(t, func() {
		storages.Get("test")
	})

	stub := StorageStub{}
	storages.storages = make(map[string]Storage)
	storages.storages["test"] = stub
	assert.Panics(t, func() {
		storages.Get("not exists")
	})

	assert.Equal(t, stub, storages.Get("test"))
}
