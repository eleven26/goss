package aliyun

import (
	"io"
	"reflect"
	"strconv"

	"github.com/eleven26/goss/core"

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
	panic("deprecated")
}

func (s *Store) Delete(key string) error {
	return s.Bucket.DeleteObject(key)
}

func (s *Store) Size(key string) (int64, error) {
	header, err := s.Bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return 0, err
	}

	length := header.Get("Content-Length")

	return strconv.ParseInt(length, 10, 64)
}

func (s *Store) Exists(key string) (bool, error) {
	return s.Bucket.IsObjectExist(key)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	return core.NewFileIterator(oss.Prefix(dir), &Chunks{prefix: dir, bucket: s.Bucket})
}

type Chunks struct {
	prefix string
	bucket *oss.Bucket
}

func (c *Chunks) Chunk(marker interface{}) (core.ListObjectResult, error) {
	var result oss.ListObjectsResult
	var err error

	if reflect.TypeOf(marker).String() == "string" {
		result, err = c.bucket.ListObjects(oss.Marker(marker.(string)), oss.Prefix(c.prefix))
	} else {
		result, err = c.bucket.ListObjects(marker.(oss.Option))
	}

	return &ListObjectResult{ossResult: result}, err
}
