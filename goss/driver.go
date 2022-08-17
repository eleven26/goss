package goss

import (
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/drivers/aliyun"
	"github.com/eleven26/goss/drivers/qiniu"
	"github.com/eleven26/goss/drivers/tencent"

	"github.com/spf13/viper"
)

const (
	Aliyun  = "aliyun"
	Tencent = "tencent"
	Qiniu   = "qiniu"
)

func defaultDriver() core.Driver {
	driver := viper.GetString("driver")

	switch driver {
	case Aliyun:
		return aliyun.NewDriver()
	case Tencent:
		return tencent.NewDriver()
	case Qiniu:
		return qiniu.NewDriver()
	default:
		panic("no default driver set.")
	}
}
