package goss

import (
	"errors"
)

type Config struct {
	Endpoint          string `yaml:"endpoint"`
	AccessKey         string `yaml:"access_key"`
	SecretKey         string `yaml:"secret_key"`
	Region            string `yaml:"region"`
	Bucket            string `yaml:"bucket"`
	UseSsl            *bool  `yaml:"use_ssl"`
	HostnameImmutable *bool  `yaml:"hostname_immutable"`
}

func (c *Config) validate() error {


	return nil
}

func (c *Config) url() string {
	prefix := "https://"
	if c.UseSsl != nil && !*c.UseSsl {
		prefix = "http://"
	}

	return prefix + c.Endpoint
}

func (c *Config) hostnameImmutable() bool {
	if c.HostnameImmutable != nil {
		return *c.HostnameImmutable
	}

	return false
}
