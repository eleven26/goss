package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	creds := credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	awsConfig := &aws.Config{
		Region:      aws.String(conf.Region),
		Endpoint:    aws.String(conf.Endpoint),
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)

	store := Store{
		s3:     svc,
		config: *conf,
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "s3"
}
