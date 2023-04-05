package goss

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Storage defines a unified interface for reading and writing cloud storage objects.
type Storage interface {
	// Put saves the content read from r to the key of oss.
	Put(key string, r io.Reader) error

	// PutFromFile saves the file pointed to by the `localPath` to the oss key.
	PutFromFile(key string, localPath string) error

	// Get gets the file pointed to by key.
	Get(key string) (io.ReadCloser, error)

	// GetString gets the file pointed to by key and returns a string.
	GetString(key string) (string, error)

	// GetBytes gets the file pointed to by key and returns a byte array.
	GetBytes(key string) ([]byte, error)

	// GetToFile saves the file pointed to by key to the localPath.
	GetToFile(key string, localPath string) error

	// Delete the file pointed to by key.
	Delete(key string) error

	// Exists determines whether the file exists.
	Exists(key string) (bool, error)

	// Files list all files in the specified directory.
	Files(dir string) ([]File, error)

	// Size fet the file size.
	Size(key string) (int64, error)

	// Store is an instance for calling APIs of different cloud storage service providers.
	Store() interface{}
}

func newStorage(conf *Config) (Storage, error) {
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

	store := Store{
		s3:     client,
		Bucket: conf.Bucket,
	}

	return &store, nil
}
