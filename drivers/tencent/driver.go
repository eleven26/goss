package tencent

import (
	"net/http"
	"net/url"

	"github.com/eleven26/goss/v2/core"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
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

	if conf.Url == "" || conf.SecretId == "" || conf.SecretKey == "" {
		return nil, core.ErrorConfigEmpty
	}

	u, err := url.Parse(conf.Url)
	if err != nil {
		return nil, err
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})

	store := Store{
		client: client,
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "tencent"
}
