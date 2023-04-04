package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DriverStub struct{}

func (d DriverStub) Storage() (Storage, error) {
	return &StorageStub{}, nil
}

func (d DriverStub) Name() string {
	return "stub"
}

func TestNew(t *testing.T) {
	kernel := New()

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
	err := kernel.UseDriver(driver)
	assert.NotNil(t, err)

	storage, err := driver.Storage()
	if err != nil {
		t.Fatal(err)
	}
	kernel.storages.storages["stub"] = storage
	err = kernel.UseDriver(driver)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, kernel.Storage, storage)
}

func TestRegisterDriver(t *testing.T) {
	kernel := Kernel{
		storages: Storages{
			storages: make(map[string]Storage),
		},
	}

	driver := DriverStub{}
	err := kernel.RegisterDriver(&driver)
	assert.Nil(t, err)

	storage, err := driver.Storage()
	assert.Nil(t, err)

	assert.Len(t, kernel.storages.storages, 1)

	s, err := kernel.storages.Get(driver.Name())
	assert.Nil(t, err)
	assert.Equal(t, s, storage)

	err = kernel.UseDriver(driver)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, kernel.Storage, storage)
}
