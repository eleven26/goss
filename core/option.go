package core

import (
	"reflect"

	"github.com/spf13/viper"
)

type Option func(Driver)

func WithViper(viper *viper.Viper) Option {
	return func(driver Driver) {
		v := reflect.ValueOf(driver)

		field := v.Elem().FieldByName("Viper")
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(viper))
		}
	}
}
