//go:build integration

package goss

import (
	"reflect"
	"testing"

	"github.com/eleven26/goss/config"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	configPath := config.UserHomeConfigPath()

	var goss Goss

	viper.Set("driver", Aliyun)
	goss = New(configPath)
	assert.Equal(t, "aliyun.Storage", reflect.TypeOf(goss.Storage.Storage()).Elem().String())

	viper.Set("driver", Tencent)
	goss = New(configPath)
	assert.Equal(t, "tencent.Storage", reflect.TypeOf(goss.Storage.Storage()).Elem().String())

	viper.Set("driver", Qiniu)
	goss = New(configPath)
	assert.Equal(t, "qiniu.Storage", reflect.TypeOf(goss.Storage.Storage()).Elem().String())

	viper.Set("driver", "not_exists")
	assert.Panics(t, func() {
		goss = New(configPath)
	})
}
