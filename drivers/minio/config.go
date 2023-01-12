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
		Endpoint:  viper.GetString("minio.endpoint"),
		AccessKey: viper.GetString("minio.access_key"),
		SecretKey: viper.GetString("minio.secret_key"),
		UseSSL:    viper.GetBool("minio.use_ssl"),
		Bucket:    viper.GetString("minio.bucket"),
	}
}
