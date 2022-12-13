package minio

import "github.com/spf13/viper"

type config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

func getConfig() *config {
	return &config{
		Endpoint:  viper.GetString("s3.endpoint"),
		AccessKey: viper.GetString("s3.access_key"),
		SecretKey: viper.GetString("s3.secret_key"),
		UseSSL:    viper.GetBool("s3.use_ssl"),
		Bucket:    viper.GetString("s3.bucket"),
	}
}
