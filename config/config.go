package config

import (
	"fmt"
	"os"

	fs "github.com/eleven26/go-filesystem"
	"github.com/spf13/viper"
)

func ReadInConfig(path string) error {
	exist, err := fs.Exists(path)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("配置文件不存在：%s\n", path)
	}

	viper.SetConfigFile(path)

	return viper.ReadInConfig()
}

func ReadInUserHomeConfig() error {
	path, err := UserHomeConfigPath()
	if err != nil {
		return err
	}

	return ReadInConfig(path)
}

func UserHomeConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file := home + "/.goss.yml"
	exist, err := fs.Exists(file)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("配置文件不存在：%s\n", file)
	}

	return file, nil
}
