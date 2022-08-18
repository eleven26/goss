package qiniu

import (
	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig()

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

	return &Storage{store: store}, nil
}

func (d Driver) Name() string {
	return "qiniu"
}
