//go:build integration

package goss

import (
	"reflect"
	"testing"

	"github.com/eleven26/goss/v2/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var goss *Goss
	var err error
	var configPath string

	configPath, err = config.UserHomeConfigPath()
	assert.Nil(t, err)

	goss, err = New(configPath)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())
}

func TestNewWithViper(t *testing.T) {
	var goss *Goss
	var err error
	var configPath string

	configPath, err = config.UserHomeConfigPath()
	assert.Nil(t, err)

	v, err := config.ReadInConfig(configPath)
	assert.Nil(t, err)

	goss, err = NewWithViper(v)
	assert.Nil(t, err)
	assert.Equal(t, "core.storage", reflect.TypeOf(goss.Storage).Elem().String())
}
