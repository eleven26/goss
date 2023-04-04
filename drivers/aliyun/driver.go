package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/v2/core"
	"github.com/spf13/viper"
)

type Driver struct {
	Viper *viper.Viper
}

func NewDriver(opts ...core.Option) core.Driver {
	driver := &Driver{}

	for _, option := range opts {
		option(driver)
	}

	return driver
}

func (d *Driver) Storage() (core.Storage, error) {
	bucket, err := d.ossBucket()
	if err != nil {
		return nil, err
	}

	store := Store{
		Bucket: bucket,
	}

	return core.NewStorage(&store), nil
}

func (d *Driver) ossBucket() (*oss.Bucket, error) {
	conf := getConfig(d.Viper)

	if conf.Endpoint == "" || conf.AccessKeyID == "" || conf.AccessKeySecret == "" {
		return nil, core.ErrorConfigEmpty
	}

	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return client.Bucket(conf.Bucket)
}

func (d Driver) Name() string {
	return "aliyun"
}
