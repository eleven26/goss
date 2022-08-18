package aliyun

import (
	"io"
	"net/http"

	"github.com/eleven26/goss/v2/core"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Store struct {
	Bucket *oss.Bucket
}

func (s *Store) Put(key string, r io.Reader) error {
	return s.Bucket.PutObject(key, r)
}

func (s *Store) PutFromFile(key string, localPath string) error {
	return s.Bucket.PutObjectFromFile(key, localPath)
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	return s.Bucket.GetObject(key)
}

func (s *Store) SaveToFile(key string, localPath string) error {
	return s.Bucket.GetObjectToFile(key, localPath)
}

func (s *Store) Delete(key string) error {
	return s.Bucket.DeleteObject(key)
}

func (s *Store) Meta(key string) (http.Header, error) {
	return s.Bucket.GetObjectDetailedMeta(key)
}

func (s *Store) Exists(key string) (bool, error) {
	return s.Bucket.IsObjectExist(key)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	chunk := func(marker oss.Option) (oss.ListObjectsResult, error) {
		return s.Bucket.ListObjects(marker)
	}

	it := newFileIterator(dir, chunk)

	return &it
}
