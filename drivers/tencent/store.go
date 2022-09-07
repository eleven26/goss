package tencent

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Store struct {
	client *cos.Client
}

func (s *Store) Put(key string, r io.Reader) error {
	response, err := s.client.Object.Put(context.Background(), key, r, nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (s *Store) PutFromFile(key string, localPath string) error {
	response, err := s.client.Object.PutFromFile(context.Background(), key, localPath, nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	resp, err := s.client.Object.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (s *Store) Delete(key string) error {
	_, err := s.client.Object.Delete(context.Background(), key)

	return err
}

func (s *Store) Size(key string) (int64, error) {
	resp, err := s.client.Object.Head(context.Background(), key, nil)
	if err != nil {
		return 0, err
	}

	return resp.ContentLength, nil
}

func (s *Store) Exists(key string) (bool, error) {
	return s.client.Object.IsExist(context.Background(), key)
}

func (s *Store) Iterator(prefix string) core.FileIterator {
	return core.NewFileIterator(&Chunks{
		prefix: prefix,
		bucket: s.client.Bucket,
	})
}

func httpError(response *cos.Response) error {
	bytes, err := io.ReadAll(response.Body)
	defer func() {
		err = response.Body.Close()
	}()
	if err != nil {
		return err
	}

	return errors.New(string(bytes))
}

type Chunks struct {
	count      int
	prefix     string
	nextMarker string
	bucket     *cos.BucketService
}

func (c *Chunks) Chunk() (core.ListObjectResult, error) {
	var opt cos.BucketGetOptions
	var result *cos.BucketGetResult
	var err error

	// 参考文档：https://cloud.tencent.com/document/product/436/7734
	// 单次返回最大的条目数量，默认值为1000，最大为1000
	// BucketGetOptions.MaxKeys 可以设置单次获取的条目数量
	if c.count == 0 {
		opt = cos.BucketGetOptions{Prefix: c.prefix}
	} else {
		opt = cos.BucketGetOptions{Prefix: c.prefix, Marker: c.nextMarker}
	}

	result, _, err = c.bucket.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}

	c.count++
	c.nextMarker = result.NextMarker

	return &ListObjectResult{result: result}, nil
}
