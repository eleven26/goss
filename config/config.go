package config

import (
	"fmt"
	"os"

	fs "github.com/eleven26/go-filesystem"
	"github.com/spf13/viper"
)

func ReadInConfig(path string) {
	exist, err := fs.Exists(path)
	if err != nil {
		panic(err)
	}
	if !exist {
		panic(fmt.Errorf("配置文件不存在：%s\n", path))
	}

	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func ReadInUserHomeConfig() {
	path := UserHomeConfigPath()
	ReadInConfig(path)
}

func UserHomeConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	file := home + "/.goss.yml"
	exist, err := fs.Exists(file)
	if err != nil {
		panic(err)
	}
	if !exist {
		panic(fmt.Errorf("配置文件不存在：%s\n", file))
	}

	return file
}
