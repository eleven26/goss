package cmd

import (
	"log"

	"github.com/eleven26/goss/goss"
)

var app goss.Goss

func Run() {
	var err error
	app, err = goss.NewFromUserHomeConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	err = Execute()
	if err != nil {
		log.Fatal(err)
	}
}
