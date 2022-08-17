package core

type Driver interface {
	Storage() Storage
	Name() string
}
