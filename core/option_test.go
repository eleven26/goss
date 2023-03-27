package core

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var _ Driver = &mockDriver{}

type mockDriver struct {
	Viper *viper.Viper
}

func (m *mockDriver) Storage() (Storage, error) {
	panic("implement me")
}

func (m *mockDriver) Name() string {
	panic("implement me")
}

func TestWithViper(t *testing.T) {
	driver := &mockDriver{}

	v := viper.New()
	v.Set("a", "b")

	opt := WithViper(v)
	opt(driver)

	assert.Equal(t, "b", driver.Viper.GetString("a"))
	assert.Equal(t, v, driver.Viper)
}
