package cmd

import (
	"github.com/eleven26/goss/config"
	"github.com/eleven26/goss/goss"
)

var app goss.Goss

func Run() {
	app = goss.New(config.UserHomeConfigPath())

	err := Execute()
	if err != nil {
		panic(err)
	}
}
