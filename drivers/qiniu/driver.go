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

func (d *Driver) Storage() core.Storage {
	config := getConfig()

	if config.Bucket == "" || config.AccessKey == "" || config.SecretKey == "" {
		panic("配置不正确")
	}

	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	store := Store{
		config:        *config,
		bucketManager: storage.NewBucketManager(mac, &cfg),
	}

	return &Storage{store: store}
}

func (d Driver) Name() string {
	return "qiniu"
}
