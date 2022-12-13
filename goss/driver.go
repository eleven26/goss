package goss

import (
	"errors"

	"github.com/eleven26/goss/drivers/minio"

	"github.com/eleven26/goss/drivers/s3"

	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/drivers/aliyun"
	"github.com/eleven26/goss/drivers/huawei"
	"github.com/eleven26/goss/drivers/qiniu"
	"github.com/eleven26/goss/drivers/tencent"

	"github.com/spf13/viper"
)

const (
	Aliyun  = "aliyun"
	Tencent = "tencent"
	Qiniu   = "qiniu"
	Huawei  = "huawei"
	S3      = "s3"
	Minio   = "minio"
)

var (
	// ErrNoDefaultDriver no default driver configured error.
	ErrNoDefaultDriver = errors.New("no default driver set")

	// ErrDriverNotExists driver not registered error.
	ErrDriverNotExists = errors.New("driver not exists")
)

// defaultDriver get the driver specified by "driver" in the configuration file.
func defaultDriver() (core.Driver, error) {
	if !viper.IsSet("driver") {
		return nil, ErrNoDefaultDriver
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
	case S3:
		return s3.NewDriver(), nil
	case Minio:
		return minio.NewDriver(), nil
	default:
		return nil, ErrDriverNotExists
	}
}
