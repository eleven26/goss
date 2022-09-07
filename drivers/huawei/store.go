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

func (s *Store) Iterator(dir string) core.FileIterator {
	return core.NewFileIterator(&Chunks{
		bucket: s.config.Bucket,
		client: s.client,
		prefix: dir,
	})
}

type Chunks struct {
	prefix     string
	nextMarker string
	bucket     string
	client     *obs.ObsClient
}

func (c *Chunks) Chunk() (core.ListObjectResult, error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = c.bucket
	input.Marker = c.nextMarker
	input.Prefix = c.prefix

	// 参考文档：https://support.huaweicloud.com/sdk-android-devg-obs/obs_26_0603.html
	// input.maxKeys 列举对象的最大数目，取值范围为1~1000，当超出范围时，按照默认的1000进行处理。
	output, err := c.client.ListObjects(input)
	if err != nil {
		return nil, err
	}

	if output.IsTruncated {
		c.nextMarker = output.NextMarker
	}

	return NewListObjectResult(output), nil
}
