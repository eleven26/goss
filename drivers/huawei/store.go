package huawei

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type Store struct {
	config
	client *obs.ObsClient
}

func (s *Store) Put(key string, r io.Reader) error {
	input := &obs.PutObjectInput{}
	input.Bucket = s.config.Bucket
	input.Key = key
	input.Body = r

	_, err := s.client.PutObject(input)

	return err
}

func (s *Store) PutFromFile(key string, localPath string) error {
	input := &obs.PutFileInput{}
	input.Bucket = s.config.Bucket
	input.Key = key
	input.SourceFile = localPath

	_, err := s.client.PutFile(input)

	return err
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	output, err := s.client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return output.Body, err
}

func (s *Store) SaveToFile(key string, localPath string) error {
	rc, err := s.Get(key)
	if err != nil {
		return err
	}

	defer func(rc io.ReadCloser) {
		err = rc.Close()
	}(rc)

	// 保存到文件 localPath
	f, _ := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	_, err = io.Copy(f, rc)

	return err
}

func (s *Store) Delete(key string) error {
	input := &obs.DeleteObjectInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	_, err := s.client.DeleteObject(input)

	return err
}

func (s *Store) Meta(key string) (http.Header, error) {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	output, err := s.client.GetObjectMetadata(input)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Set("Content-Length", strconv.FormatInt(output.ContentLength, 10))

	return header, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.Meta(key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s Store) ListObjects(marker oss.Option) (oss.ListObjectsResult, error) {
	panic("")
}

func (s *Store) Iterator(dir string) core.FileIterator {
	panic("")
}
