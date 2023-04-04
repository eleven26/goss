package goss

import (
	"fmt"

	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/drivers/aliyun"
	"github.com/eleven26/goss/drivers/huawei"
	"github.com/eleven26/goss/drivers/minio"
	"github.com/eleven26/goss/drivers/qiniu"
	"github.com/eleven26/goss/drivers/s3"
	"github.com/eleven26/goss/drivers/tencent"
)

const (
	Aliyun  = "aliyun"
	Tencent = "tencent"
	Qiniu   = "qiniu"
	Huawei  = "huawei"
	S3      = "s3"
	Minio   = "minio"
)

// defaultDriver get the driver specified by "driver" in the configuration file.
func defaultDriver(driver string, opts ...core.Option) (core.Driver, error) {
	switch driver {
	case Aliyun:
		return aliyun.NewDriver(opts...), nil
	case Tencent:
		return tencent.NewDriver(opts...), nil
	case Qiniu:
		return qiniu.NewDriver(opts...), nil
	case Huawei:
		return huawei.NewDriver(opts...), nil
	case S3:
		return s3.NewDriver(opts...), nil
	case Minio:
		return minio.NewDriver(opts...), nil
	default:
		return nil, fmt.Errorf("driver not exists: %s", driver)
	}
}
