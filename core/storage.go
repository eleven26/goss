package core

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type Storage interface {
	// Put 将从 r 读取的内容保存到 oss 的 key
	Put(key string, r io.Reader) error
	// PutFromFile 将本地路径 localPath 指向的文件保存到 oss 的 key
	PutFromFile(key string, localPath string) error

	// Get 获取 key 指向的文件
	Get(key string) (io.ReadCloser, error)
	// GetString 获取 key 指向的文件，返回字符串
	GetString(key string) (string, error)
	// GetBytes 获取 key 指向的文件，返回字节数组
	GetBytes(key string) ([]byte, error)
	// GetToFile 保存 key 指向的文件到本地 localPath
	GetToFile(key string, localPath string) error

	// Delete 删除 key 指向的文件
	Delete(key string) error
	// Exists 判断文件是否存在
	Exists(key string) (bool, error)
	// Files 列出指定目录下的所有文件
	Files(dir string) ([]File, error)
	// Size 获取文件大小
	Size(key string) (int64, error)

	Store() interface{}
}

type storage struct {
	store Store
}

func NewStorage(store Store) Storage {
	return &storage{
		store: store,
	}
}

func (s *storage) Put(key string, r io.Reader) error {
	return s.store.Put(key, r)
}

func (s *storage) PutFromFile(key string, localPath string) error {
	return s.store.PutFromFile(key, localPath)
}

func (s *storage) Get(key string) (reader io.ReadCloser, err error) {
	return s.store.Get(key)
}

func (s *storage) GetString(key string) (string, error) {
	bs, err := s.GetBytes(key)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (s *storage) GetBytes(key string) (bytes []byte, err error) {
	rc, err := s.store.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	return ioutil.ReadAll(rc)
}

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

func (s *storage) Size(key string) (int64, error) {
	header, err := s.store.Meta(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(header.Get("Content-Length"), 10, 64)
}

func (s *storage) Delete(key string) error {
	return s.store.Delete(key)
}

func (s *storage) Exists(key string) (bool, error) {
	return s.store.Exists(key)
}

func (s *storage) Files(dir string) ([]File, error) {
	return s.store.Iterator(dir).All()
}

func (s *storage) Store() interface{} {
	return s.store
}
