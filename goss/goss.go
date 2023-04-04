package goss

import (
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/internal/config"
	"github.com/spf13/viper"
)

// Goss is the wrapper for core.Kernel
type Goss struct {
	core.Kernel
}

// New creates a new instance based on the configuration file pointed to by configPath.
func New(configPath string) (*Goss, error) {
	v, err := config.ReadInConfig(configPath)
	if err != nil {
		return nil, err
	}

	goss := Goss{
		core.New(),
	}

	driver, err := defaultDriver(v.GetString("driver"), core.WithViper(v))
	if err != nil {
		return nil, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return nil, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return nil, err
	}

	return &goss, nil
}

// NewWithViper creates a new instance based on the configuration file pointed to by viper.
func NewWithViper(v *viper.Viper) (*Goss, error) {
	goss := Goss{
		core.New(),
	}

	driver, err := defaultDriver(v.GetString("driver"), core.WithViper(v))
	if err != nil {
		return nil, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return nil, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return nil, err
	}

	return &goss, nil
}

// NewFromUserHomeConfigPath creates a new instance based on the configuration file pointed to by user home directory.
func NewFromUserHomeConfigPath() (*Goss, error) {
	path, err := config.UserHomeConfigPath()
	if err != nil {
		return nil, err
	}

	return New(path)
}
