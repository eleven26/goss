package qiniu

import (
	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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
	conf := getConfig(d.Viper)

	if conf.Bucket == "" || conf.AccessKey == "" || conf.SecretKey == "" {
		return nil, core.ErrorConfigEmpty
	}

	mac := qbox.NewMac(conf.AccessKey, conf.SecretKey)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	store := Store{
		config:        *conf,
		bucketManager: storage.NewBucketManager(mac, &cfg),
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "qiniu"
}
