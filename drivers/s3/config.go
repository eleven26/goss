package s3

import "github.com/spf13/viper"

type config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

func getConfig() *config {
	return &config{
		Endpoint:  viper.GetString("s3.endpoint"),
		AccessKey: viper.GetString("s3.access_key"),
		SecretKey: viper.GetString("s3.secret_key"),
		Region:    viper.GetString("s3.region"),
		Bucket:    viper.GetString("s3.bucket"),
	}
}
