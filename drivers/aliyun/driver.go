package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() core.Storage {
	config := getConfig()

	if config.Endpoint == "" || config.AccessKeyID == "" || config.AccessKeySecret == "" {
		panic("配置不正确")
	}

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		panic(err)
	}

	store := Store{
		Bucket: bucket,
	}

	return &Storage{store: store}
}

func (d Driver) Name() string {
	return "aliyun"
}
