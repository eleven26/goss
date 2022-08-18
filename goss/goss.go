package goss

import (
	"github.com/eleven26/goss/v2/config"
	"github.com/eleven26/goss/v2/core"
	"github.com/eleven26/goss/v2/drivers/aliyun"
	"github.com/eleven26/goss/v2/drivers/qiniu"
	"github.com/eleven26/goss/v2/drivers/tencent"
)

type Goss struct {
	core.Kernel
}

func New(configPath string) (Goss, error) {
	err := config.ReadInConfig(configPath)
	if err != nil {
		return Goss{}, err
	}

	goss := Goss{
		core.New(),
	}

	driver, err := defaultDriver()
	if err != nil {
		return Goss{}, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return Goss{}, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return Goss{}, err
	}

	return goss, nil
}

func (g *Goss) RegisterAliyunDriver() error {
	return g.RegisterDriver(aliyun.NewDriver())
}

func (g *Goss) RegisterTencentDriver() error {
	return g.RegisterDriver(tencent.NewDriver())
}

func (g *Goss) RegisterQiniuDriver() error {
	return g.RegisterDriver(qiniu.NewDriver())
}

func NewFromUserHomeConfigPath() (Goss, error) {
	path, err := config.UserHomeConfigPath()
	if err != nil {
		return Goss{}, err
	}

	return New(path)
}
