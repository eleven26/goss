package goss

import (
	"errors"

	"github.com/eleven26/goss/drivers/huawei"

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
	Huawei  = "huawei"
)

var (
	errorNoDefaultDriver = errors.New("no default driver set")
	errorDriverNotExists = errors.New("driver not exists")
)

func defaultDriver() (core.Driver, error) {
	if !viper.IsSet("driver") {
		return nil, errorNoDefaultDriver
	}

	driver := viper.GetString("driver")

	switch driver {
	case Aliyun:
		return aliyun.NewDriver(), nil
	case Tencent:
		return tencent.NewDriver(), nil
	case Qiniu:
		return qiniu.NewDriver(), nil
	case Huawei:
		return huawei.NewDriver(), nil
	default:
		return nil, errorDriverNotExists
	}
}
