package tencent

import (
	"net/http"
	"net/url"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() core.Storage {
	config := getConfig()

	if config.Url == "" || config.SecretId == "" || config.SecretKey == "" {
		panic("配置不正确")
	}

	u, err := url.Parse(config.Url)
	if err != nil {
		panic(err)
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretId,
			SecretKey: config.SecretKey,
		},
	})

	store := Store{
		client: client,
	}

	return &Storage{store: store}
}

func (d Driver) Name() string {
	return "tencent"
}
