package aliyun

import "github.com/spf13/viper"

type config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	AccessKeySecret string
}

func getConfig(viper *viper.Viper) *config {
	return &config{
		Endpoint:        viper.GetString("aliyun.endpoint"),
		Bucket:          viper.GetString("aliyun.bucket"),
		AccessKeyID:     viper.GetString("aliyun.access_key_id"),
		AccessKeySecret: viper.GetString("aliyun.access_key_secret"),
	}
}
