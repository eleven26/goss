package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DriverStub struct{}

func (d DriverStub) Storage() Storage {
	return &StorageStub{}
}

func (d DriverStub) Name() string {
	return "stub"
}

func TestNew(t *testing.T) {
	kernel := New()

	assert.Equal(t, kernel.driver, "")
	assert.Nil(t, kernel.Storage)
	assert.NotNil(t, kernel.storages)
}

func TestUseDriver(t *testing.T) {
	kernel := Kernel{
		storages: Storages{
			storages: make(map[string]Storage),
		},
	}

	driver := DriverStub{}
	kernel.storages.storages["stub"] = driver.Storage()
	kernel.UseDriver(driver)

	assert.Equal(t, kernel.Storage, driver.Storage())
}

func TestRegisterDriver(t *testing.T) {
	kernel := Kernel{
		storages: Storages{
			storages: make(map[string]Storage),
		},
	}

	driver := DriverStub{}
	kernel.RegisterDriver(&driver)

	assert.Len(t, kernel.storages.storages, 1)
	assert.Equal(t, kernel.storages.Get(driver.Name()), driver.Storage())

	kernel.UseDriver(driver)
	assert.Equal(t, kernel.Storage, driver.Storage())
}
