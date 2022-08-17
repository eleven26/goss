package tencent

import (
	"context"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Store struct {
	client *cos.Client
}

func (s *Store) Put(key string, localPath string) (*cos.Response, error) {
	return s.client.Object.PutFromFile(context.Background(), key, localPath, nil)
}

func (s *Store) Get(key string) (*cos.Response, error) {
	return s.client.Object.Get(context.Background(), key, nil)
}

func (s *Store) SaveToFile(key string, localPath string) (*cos.Response, error) {
	return s.client.Object.GetToFile(context.Background(), key, localPath, nil)
}

func (s *Store) Delete(key string) error {
	_, err := s.client.Object.Delete(context.Background(), key)

	return err
}

func (s *Store) Head(key string) (*cos.Response, error) {
	return s.client.Object.Head(context.Background(), key, nil)
}

func (s *Store) Exists(key string) (bool, error) {
	return s.client.Object.IsExist(context.Background(), key)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	chunk := func(opt *cos.BucketGetOptions) (*cos.BucketGetResult, *cos.Response, error) {
		return s.client.Bucket.Get(context.Background(), opt)
	}

	it := newFileIterator(dir)
	it.nextChunk = chunk

	return &it
}
