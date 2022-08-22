package huawei

import (
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig()

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	store := Store{
		client: client,
		config: *conf,
	}

	return core.NewStorage(&store), nil
}

func getClient() (*obs.ObsClient, error) {
	conf := getConfig()

	if conf.Endpoint == "" || conf.Location == "" || conf.Bucket == "" || conf.AccessKey == "" || conf.SecretKey == "" {
		return nil, core.ErrorConfigEmpty
	}

	return obs.New(conf.AccessKey, conf.SecretKey, conf.Endpoint)
}

func (d Driver) Name() string {
	return "huawei"
}
