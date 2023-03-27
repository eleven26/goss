package tencent

import "github.com/spf13/viper"

type config struct {
	Url       string
	SecretId  string
	SecretKey string
}

func getConfig(viper *viper.Viper) *config {
	return &config{
		Url:       viper.GetString("tencent.url"),
		SecretId:  viper.GetString("tencent.secret_id"),
		SecretKey: viper.GetString("tencent.secret_key"),
	}
}
