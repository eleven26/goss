package minio

import (
	"github.com/eleven26/goss/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig()

	if conf.Bucket == "" || conf.AccessKey == "" || conf.SecretKey == "" || conf.Endpoint == "" {
		return nil, core.ErrorConfigEmpty
	}

	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	store := Store{
		client: minioClient,
		config: *conf,
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "minio"
}
