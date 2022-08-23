package core

import "errors"

var ErrorConfigEmpty = errors.New("configuration not correct")

type Driver interface {
	Storage() (Storage, error)
	Name() string
}
