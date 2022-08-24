package core

import (
	"strings"

	"github.com/spf13/viper"
)

// Kernel is the core struct of goss, it plays the role of a driver manager.
type Kernel struct {
	driver   string
	storages Storages
	Storage  Storage
}

// New create a new instance of Kernel.
func New() Kernel {
	app := Kernel{
		driver:   strings.ToLower(viper.GetString("driver")),
		storages: Storages{},
	}

	return app
}

// UseDriver is used to switch to the specified driver.
func (a *Kernel) UseDriver(driver Driver) error {
	storage, err := a.storages.Get(strings.ToLower(driver.Name()))
	if err != nil {
		return err
	}

	a.Storage = storage

	return nil
}

// RegisterDriver is used to register new driver.
func (a *Kernel) RegisterDriver(driver Driver) error {
	storage, err := driver.Storage()
	if err != nil {
		return err
	}

	return a.storages.Register(driver.Name(), storage)
}
