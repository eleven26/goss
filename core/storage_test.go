package core

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eleven26/goss/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type StoreStub struct {
	mock.Mock
}

func (s *StoreStub) Put(key string, r io.Reader) error {
	args := s.Called(key, r)

	return args.Error(0)
}

func (s *StoreStub) PutFromFile(key string, localPath string) error {
	args := s.Called(key, localPath)

	return args.Error(0)
}

func (s *StoreStub) Get(key string) (io.ReadCloser, error) {
	args := s.Called(key)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (s *StoreStub) SaveToFile(key string, localPath string) error {
	args := s.Called(key, localPath)

	return args.Error(0)
}

func (s *StoreStub) Delete(key string) error {
	args := s.Called(key)

	return args.Error(0)
}

func (s *StoreStub) Size(key string) (int64, error) {
	args := s.Called(key)

	return args.Get(0).(int64), args.Error(1)
}

func (s *StoreStub) Exists(key string) (bool, error) {
	args := s.Called(key)

	return args.Bool(0), args.Error(1)
}

func (s *StoreStub) Iterator(marker string) FileIterator {
	args := s.Called(marker)

	return args.Get(0).(FileIterator)
}

var (
	bs = []byte("test")

	key          = "test/foo.txt"
	testdata     string
	fooPath      string
	localFooPath string
)

func init() {
	testdata = filepath.Join(utils.RootDir(), "testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")
}

func r() io.Reader {
	return strings.NewReader("test")
}

func rc() io.ReadCloser {
	return io.NopCloser(strings.NewReader("test"))
}

func TestPut(t *testing.T) {
	store := new(StoreStub)
	store.On("Put", key, r()).Return(nil)

	storage := storage{
		store: store,
	}
	err := storage.Put(key, r())
	assert.Nil(t, err)

	store.AssertExpectations(t)
}

func TestPutFromFile(t *testing.T) {
	store := new(StoreStub)
	store.On("PutFromFile", key, fooPath).Return(nil)

	storage := storage{
		store: store,
	}
	err := storage.PutFromFile(key, fooPath)
	assert.Nil(t, err)

	store.AssertExpectations(t)
}

func TestStorageGet(t *testing.T) {
	store := new(StoreStub)
	store.On("Get", key).Return(rc(), nil)

	storage := storage{
		store: store,
	}
	rc1, err := storage.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, rc1, rc())

	store.AssertExpectations(t)
}

func TestStorageGetBytes(t *testing.T) {
	store := new(StoreStub)
	store.On("Get", key).Return(rc(), nil)

	storage := storage{
		store: store,
	}
	bs1, err := storage.GetBytes(key)
	assert.Nil(t, err)
	assert.Equal(t, bs1, bs)

	store.AssertExpectations(t)
}

func TestStorageGetString(t *testing.T) {
	store := new(StoreStub)
	store.On("Get", key).Return(rc(), nil)

	storage := storage{
		store: store,
	}
	s, err := storage.GetString(key)
	assert.Nil(t, err)
	assert.Equal(t, "test", s)

	store.AssertExpectations(t)
}

func TestStorageGetToFile(t *testing.T) {
	store := new(StoreStub)
	store.On("Get", key).Return(rc(), nil)

	storage := storage{
		store: store,
	}
	err := storage.GetToFile(key, localFooPath)
	assert.Nil(t, err)

	assert.FileExists(t, localFooPath)

	bs1, err := os.ReadFile(localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, bs1, bs)

	store.AssertExpectations(t)
}

type FileIteratorStub struct {
	mock.Mock
}

func (f *FileIteratorStub) HasNext() bool {
	args := f.Called()

	return args.Bool(0)
}

func (f *FileIteratorStub) Next() (File, error) {
	args := f.Called()

	return args.Get(0).(File), args.Error(1)
}

func (f *FileIteratorStub) All() ([]File, error) {
	args := f.Called()

	return args.Get(0).([]File), args.Error(1)
}

func (f *FileIteratorStub) GetNextChunk() error {
	args := f.Called()

	return args.Error(0)
}

func TestStorageOthers(t *testing.T) {
	it := new(FileIteratorStub)
	it.On("All").Return(make([]File, 0), nil)

	store := new(StoreStub)
	var s int64 = 3
	store.On("Size", key).Return(s, nil)
	store.On("Delete", key).Return(nil)
	store.On("Exists", key).Return(false, nil)
	store.On("Iterator", key).Return(it, nil)

	storage := storage{
		store: store,
	}
	s1, err := storage.Size(key)
	assert.Nil(t, err)
	assert.Equal(t, s1, s)

	err = storage.Delete(key)
	assert.Nil(t, err)

	exists, err := storage.Exists(key)
	assert.Nil(t, err)
	assert.False(t, exists)

	files, err := storage.Files(key)
	assert.Nil(t, err)
	assert.Empty(t, files)

	store1 := storage.Store()
	assert.Equal(t, store1, store)

	store.AssertExpectations(t)
}
