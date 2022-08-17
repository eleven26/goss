package core

import (
	"strings"

	"github.com/spf13/viper"
)

type Kernel struct {
	driver   string
	storages Storages
	Storage  Storage
}

func New() Kernel {
	app := Kernel{
		driver:   strings.ToLower(viper.GetString("driver")),
		storages: Storages{},
	}

	return app
}

func (a *Kernel) UseDriver(driver Driver) {
	a.Storage = a.storages.Get(strings.ToLower(driver.Name()))
}

func (a *Kernel) RegisterDriver(driver Driver) {
	a.storages.Register(driver.Name(), driver.Storage())
}
