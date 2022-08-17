package goss

import (
	"github.com/eleven26/goss/config"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/drivers/aliyun"
	"github.com/eleven26/goss/drivers/qiniu"
	"github.com/eleven26/goss/drivers/tencent"
)

type Goss struct {
	core.Kernel
}

func New(configPath string) Goss {
	config.ReadInConfig(configPath)

	goss := Goss{
		core.New(),
	}

	goss.RegisterDriver(aliyun.NewDriver())
	goss.RegisterDriver(tencent.NewDriver())
	goss.RegisterDriver(qiniu.NewDriver())

	goss.UseDriver(defaultDriver())

	return goss
}

func Storage(configPath string) core.Storage {
	goss := New(configPath)

	return goss.Storage
}
