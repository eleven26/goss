package huawei

import (
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
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

	client, err := d.getClient()
	if err != nil {
		return nil, err
	}

	store := Store{
		client: client,
		config: *conf,
	}

	return core.NewStorage(&store), nil
}

func (d *Driver) getClient() (*obs.ObsClient, error) {
	conf := getConfig(d.Viper)

	if conf.Endpoint == "" || conf.Location == "" || conf.Bucket == "" || conf.AccessKey == "" || conf.SecretKey == "" {
		return nil, core.ErrorConfigEmpty
	}

	return obs.New(conf.AccessKey, conf.SecretKey, conf.Endpoint)
}

func (d Driver) Name() string {
	return "huawei"
}
