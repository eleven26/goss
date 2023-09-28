package goss

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// max keys for list objects
var maxKeys int32 = 1000

// Store defines interface for cloud storage.
type Store interface {
	// Put saves the content read from r to the key of oss.
	Put(ctx context.Context, key string, r io.Reader) error

	// PutFromFile saves the file pointed to by the `localPath` to the oss key.
	PutFromFile(ctx context.Context, key string, localPath string) error

	// Get gets the file pointed to by key.
	Get(ctx context.Context, key string) (io.ReadCloser, error)

	// GetString gets the file pointed to by key and returns a string.
	GetString(ctx context.Context, key string) (string, error)

	// GetBytes gets the file pointed to by key and returns a byte array.
	GetBytes(ctx context.Context, key string) ([]byte, error)

	// GetToFile saves the file pointed to by key to the localPath.
	GetToFile(ctx context.Context, key string, localPath string) error

	// Delete the file pointed to by key.
	Delete(ctx context.Context, key string) error

	// Exists determines whether the file exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Files list all files in the specified directory.
	Files(ctx context.Context, dir string) ([]File, error)

	// Size fet the file size.
	Size(ctx context.Context, key string) (int64, error)
}

func newStore(conf *Config) (Store, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			HostnameImmutable: conf.hostnameImmutable(),
			URL:               conf.url(),
			SigningRegion:     conf.Region,
		}, nil
	})
	creds := credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, "")
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(creds),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
		awsConfig.WithRegion(conf.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &store{
		s3:     client,
		Bucket: conf.Bucket,
	}, nil
}

var _ Store = &store{}

type store struct {
	s3     *s3.Client
	Bucket string
}

func (s *store) Put(ctx context.Context, key string, r io.Reader) error {
	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return s.putFile(ctx, key, bytes.NewReader(bs))
}

func (s *store) putFile(ctx context.Context, key string, f io.ReadSeeker) error {
	input := &s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.s3.PutObject(ctx, input)

	return err
}

func (s *store) PutFromFile(ctx context.Context, key string, localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}

	return s.putFile(ctx, key, f)
}

func (s *store) getObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	return s.s3.GetObject(ctx, input)
}

func (s *store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := s.getObject(ctx, key)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

func (s *store) head(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	return s.s3.HeadObject(ctx, input)
}

func (s *store) Size(ctx context.Context, key string) (int64, error) {
	output, err := s.head(ctx, key)
	if err != nil {
		return 0, err
	}

	return output.ContentLength, nil
}

func (s *store) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.head(ctx, key)
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

func (s *store) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// GetBytes gets the file pointed to by key and returns a byte array.
func (s *store) GetBytes(ctx context.Context, key string) (bytes []byte, err error) {
	rc, err := s.Get(ctx, key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	return io.ReadAll(rc)
}

// GetString gets the file pointed to by key and returns a string.
func (s *store) GetString(ctx context.Context, key string) (string, error) {
	bs, err := s.GetBytes(ctx, key)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// GetToFile saves the file pointed to by key to the localPath.
func (s *store) GetToFile(ctx context.Context, key string, localPath string) (err error) {
	rc, err := s.Get(ctx, key)
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

// Files list all files in the given prefix.
func (s *store) Files(ctx context.Context, dir string) ([]File, error) {
	var continuationToken *string
	var count int32
	var files []File

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.Bucket),
			ContinuationToken: continuationToken,
			Prefix:            aws.String(dir),
			MaxKeys:           maxKeys,
		}

		output, err := s.s3.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Contents {
			if strings.HasSuffix(*item.Key, "/") {
				continue
			}

			files = append(files, &file{item})
			count++
		}

		if !output.IsTruncated {
			break
		}

		continuationToken = output.NextContinuationToken
	}

	return files, nil
}
