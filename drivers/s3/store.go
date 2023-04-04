package s3

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/eleven26/goss/v2/core"
)

type Store struct {
	s3 *s3.S3
	config
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
	_, err := s.s3.PutObject(input)

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

	return s.s3.GetObject(input)
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

	return s.s3.HeadObject(input)
}

func (s *Store) Size(key string) (int64, error) {
	output, err := s.head(key)
	if err != nil {
		return 0, err
	}

	return *output.ContentLength, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.head(key)
	if err != nil {
		if _, ok := err.(awserr.Error); ok {
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				if reqErr.StatusCode() == http.StatusNotFound {
					return false, nil
				}
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

	_, err := s.s3.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Iterator(prefix string) core.FileIterator {
	return core.NewFileIterator(NewChunks(s.Bucket, prefix, s.s3))
}
