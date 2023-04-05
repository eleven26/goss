package goss

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

var _ Storage = &Store{}

type Store struct {
	s3     *s3.Client
	Bucket string
}

func (s *Store) Files(dir string) ([]File, error) {
	return newFileIterator(newChunks(s.Bucket, dir, s.s3)).All()
}

func (s *Store) Store() interface{} {
	return s
}

func (s *Store) Put(key string, r io.Reader) error {
	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return s.putFile(key, bytes.NewReader(bs))
}

func (s *Store) putFile(key string, f io.ReadSeeker) error {
	input := &s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.s3.PutObject(context.TODO(), input)

	return err
}

func (s *Store) PutFromFile(key string, localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}

	return s.putFile(key, f)
}

func (s *Store) getObject(key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	return s.s3.GetObject(context.TODO(), input)
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	output, err := s.getObject(key)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

func (s *Store) head(key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	return s.s3.HeadObject(context.TODO(), input)
}

func (s *Store) Size(key string) (int64, error) {
	output, err := s.head(key)
	if err != nil {
		return 0, err
	}

	return output.ContentLength, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.head(key)
	if err != nil {
		if e, ok := err.(*smithy.OperationError); ok {
			if strings.Contains(e.Err.Error(), "404") {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (s *Store) Delete(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3.DeleteObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

// GetBytes gets the file pointed to by key and returns a byte array.
func (s *Store) GetBytes(key string) (bytes []byte, err error) {
	rc, err := s.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	return io.ReadAll(rc)
}

// GetString gets the file pointed to by key and returns a string.
func (s *Store) GetString(key string) (string, error) {
	bs, err := s.GetBytes(key)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// GetToFile saves the file pointed to by key to the localPath.
func (s *Store) GetToFile(key string, localPath string) (err error) {
	rc, err := s.Get(key)
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
