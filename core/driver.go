package core

import "errors"

// ErrorConfigEmpty is configuration empty error.
var ErrorConfigEmpty = errors.New("configuration not correct")

// Driver is the abstraction for different cloud storage service provider.
type Driver interface {
	// Storage is the abstraction for cloud storage.
	Storage() (Storage, error)

	// Name is the name of provider.
	Name() string
}
