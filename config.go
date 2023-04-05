package goss

import (
	"errors"
)

type Config struct {
	Endpoint          string
	AccessKey         string
	SecretKey         string
	Region            string
	Bucket            string
	UseSsl            *bool
	HostnameImmutable *bool
}

func (c *Config) validate() error {
	if c.Bucket == "" || c.AccessKey == "" || c.SecretKey == "" || c.Endpoint == "" || c.Region == "" {
		return errors.New("configuration not correct")
	}

	return nil
}

func (c *Config) url() string {
	prefix := "http://"
	if *c.UseSsl {
		prefix = "https://"
	}

	return prefix + c.Endpoint
}

func (c *Config) hostnameImmutable() bool {
	return *c.HostnameImmutable
}
