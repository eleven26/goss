package core

import "errors"

var ErrorConfigEmpty = errors.New("配置不正确")

type Driver interface {
	Storage() (Storage, error)
	Name() string
}
