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
	var goss *Goss
	var err error
	var configPath string

	configPath, err = config.UserHomeConfigPath()
	assert.Nil(t, err)

	viper.Set("driver", Aliyun)
	goss, err = New(configPath)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())

	viper.Set("driver", Tencent)
	goss, err = New(configPath)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())

	viper.Set("driver", Qiniu)
	goss, err = New(configPath)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())

	viper.Set("driver", Huawei)
	goss, err = New(configPath)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())

	viper.Set("driver", "not_exists")
	goss, err = New(configPath)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrDriverNotExists)
}
