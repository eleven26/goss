package minio

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"

	"github.com/eleven26/goss/core"
)

type Store struct {
	client *minio.Client
	config
}

func (s *Store) Put(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(context.Background(), s.config.Bucket, key, buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})

	return err
}

func (s *Store) putFile(key string, f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	size := fi.Size()

	_, err = s.client.PutObject(context.Background(), s.config.Bucket, key, f, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})

	return err
}

func (s *Store) PutFromFile(key string, localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}

	return s.putFile(key, f)
}

func (s *Store) getObject(key string) (*minio.Object, error) {
	return s.client.GetObject(context.Background(), s.config.Bucket, key, minio.GetObjectOptions{})
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	object, err := s.getObject(key)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (s *Store) stat(key string) (minio.ObjectInfo, error) {
	return s.client.StatObject(context.Background(), s.config.Bucket, key, minio.StatObjectOptions{})
}

func (s *Store) Size(key string) (int64, error) {
	info, err := s.stat(key)
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.stat(key)
	if err != nil {
		if errResponse, ok := err.(minio.ErrorResponse); ok {
			if errResponse.StatusCode == http.StatusNotFound {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (s *Store) Delete(key string) error {
	return s.client.RemoveObject(context.Background(), s.config.Bucket, key, minio.RemoveObjectOptions{})
}

func (s *Store) Iterator(prefix string) core.FileIterator {
	return core.NewFileIterator(NewChunks(s.Bucket, prefix, s.client))
}
