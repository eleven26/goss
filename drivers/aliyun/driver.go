package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig()

	if conf.Endpoint == "" || conf.AccessKeyID == "" || conf.AccessKeySecret == "" {
		return nil, core.ErrorConfigEmpty
	}

	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(conf.Bucket)
	if err != nil {
		return nil, err
	}

	store := Store{
		Bucket: bucket,
	}

	return &Storage{store: store}, nil
}

func (d Driver) Name() string {
	return "aliyun"
}
