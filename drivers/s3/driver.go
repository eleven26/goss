package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eleven26/goss/v2/core"
	"github.com/spf13/viper"
)

type Driver struct {
	Viper *viper.Viper
}

func NewDriver(opts ...core.Option) core.Driver {
	driver := &Driver{}

	for _, option := range opts {
		option(driver)
	}

	return driver
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig(d.Viper)

	if conf.Bucket == "" || conf.AccessKey == "" || conf.SecretKey == "" || conf.Endpoint == "" || conf.Region == "" {
		return nil, core.ErrorConfigEmpty
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "oss",
			URL:           "https://" + conf.Endpoint,
			SigningRegion: conf.Region,
		}, nil
	})
	creds := credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, "")
	cfg, err := config2.LoadDefaultConfig(context.TODO(), config2.WithCredentialsProvider(creds), config2.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	store := Store{
		s3:     client,
		config: *conf,
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "s3"
}
