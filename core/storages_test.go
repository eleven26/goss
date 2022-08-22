package core

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StorageStub struct{}

func (s StorageStub) Put(key string, r io.Reader) error {
	panic("implement me")
}

func (s StorageStub) PutFromFile(key string, localPath string) error {
	panic("implement me")
}

func (s StorageStub) Get(key string) (io.ReadCloser, error) {
	panic("implement me")
}

func (s StorageStub) GetString(key string) (string, error) {
	panic("implement me")
}

func (s StorageStub) GetBytes(key string) ([]byte, error) {
	panic("implement me")
}

func (s StorageStub) Delete(key string) error {
	panic("implement me")
}

func (s StorageStub) GetToFile(key string, localPath string) error {
	panic("implement me")
}

func (s StorageStub) Exists(key string) (bool, error) {
	panic("implement me")
}

func (s StorageStub) Files(dir string) ([]File, error) {
	panic("implement me")
}

func (s StorageStub) Size(key string) (int64, error) {
	panic("implement me")
}

func (s StorageStub) Store() interface{} {
	panic("implement me")
}

func TestRegister(t *testing.T) {
	storages := Storages{}

	assert.Nil(t, storages.storages)
	stub := StorageStub{}
	err := storages.Register("test", stub)
	assert.Nil(t, err)
	assert.NotNil(t, storages.storages)

	assert.Len(t, storages.storages, 1)

	s, ok := storages.storages["test"]
	assert.True(t, ok)
	assert.Equal(t, s, stub)

	err = storages.Register("test", stub)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, errExistsDriver)

	err = storages.Register("test1", StorageStub{})
	assert.Nil(t, err)
	assert.Len(t, storages.storages, 2)
}

func TestGet(t *testing.T) {
	storages := Storages{}

	s, err := storages.Get("test")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	stub := StorageStub{}
	storages.storages = make(map[string]Storage)
	storages.storages["test"] = stub

	s, err = storages.Get("not exists")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = storages.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, stub, s)
}
