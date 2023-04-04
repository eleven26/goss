package config

import (
	"fmt"
	"os"
	"path/filepath"

	fs "github.com/eleven26/go-filesystem"
	"github.com/spf13/viper"
)

// File The name of configuration file.
const File = ".goss.yml"

// ReadInConfig Load the configuration file from the specified path.
func ReadInConfig(path string) (*viper.Viper, error) {
	exist, err := fs.Exists(path)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("Configuration file not exists：%s\n", path)
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// ReadInUserHomeConfig Load the configuration file from user home directory.
func ReadInUserHomeConfig() (*viper.Viper, error) {
	path, err := UserHomeConfigPath()
	if err != nil {
		return nil, err
	}

	return ReadInConfig(path)
}

// UserHomeConfigPath Get the user home directory configuration path.
func UserHomeConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file := filepath.Join(home, File)
	exist, err := fs.Exists(file)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("Configuration file not exists：%s\n", file)
	}

	return file, nil
}
