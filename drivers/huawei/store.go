package huawei

import (
	"io"

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
	panic("deprecated")
}

func (s *Store) Delete(key string) error {
	input := &obs.DeleteObjectInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	_, err := s.client.DeleteObject(input)

	return err
}

func (s *Store) meta(key string) (*obs.GetObjectMetadataOutput, error) {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	output, err := s.client.GetObjectMetadata(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *Store) Size(key string) (int64, error) {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = s.config.Bucket
	input.Key = key

	output, err := s.client.GetObjectMetadata(input)
	if err != nil {
		return 0, err
	}

	return output.ContentLength, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.meta(key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s Store) listObjects(input *obs.ListObjectsInput) (output *obs.ListObjectsOutput, err error) {
	return s.client.ListObjects(input)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	chunk := func(marker interface{}) (core.ListObjectResult, error) {
		input := &obs.ListObjectsInput{}
		input.Bucket = s.config.Bucket
		input.Marker = marker.(string)

		output, err := s.listObjects(input)
		if err != nil {
			return nil, err
		}

		return &ListObjectResult{
			output: output,
		}, nil
	}

	return core.NewFileIterator(dir, chunk)
}
