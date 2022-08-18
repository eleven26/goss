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

func (a *Kernel) UseDriver(driver Driver) error {
	storage, err := a.storages.Get(strings.ToLower(driver.Name()))
	if err != nil {
		return err
	}

	a.Storage = storage

	return nil
}

func (a *Kernel) RegisterDriver(driver Driver) error {
	storage, err := driver.Storage()
	if err != nil {
		return err
	}

	return a.storages.Register(driver.Name(), storage)
}
