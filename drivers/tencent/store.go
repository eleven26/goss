package tencent

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
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

func (s *Store) SaveToFile(key string, localPath string) error {
	panic("deprecated")
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

func (s Store) getWithOpt(opt *cos.BucketGetOptions) (*cos.BucketGetResult, *cos.Response, error) {
	return s.client.Bucket.Get(context.Background(), opt)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	chunk := func(marker interface{}) (core.ListObjectResult, error) {
		var result *cos.BucketGetResult
		var err error

		if opt, ok := marker.(*cos.BucketGetOptions); ok {
			result, _, err = s.getWithOpt(opt)
		} else {
			opt = &cos.BucketGetOptions{Prefix: dir}
			result, _, err = s.getWithOpt(opt)
		}

		if err != nil {
			return nil, err
		}

		return &ListObjectResult{result: result}, nil
	}

	return core.NewFileIterator(dir, chunk)
}

func httpError(response *cos.Response) error {
	bytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err = response.Body.Close()
	}()
	if err != nil {
		return err
	}

	return errors.New(string(bytes))
}
