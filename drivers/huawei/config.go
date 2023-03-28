package huawei

import "github.com/spf13/viper"

type config struct {
	Endpoint  string
	Location  string
	Bucket    string
	AccessKey string
	SecretKey string
}

func getConfig(viper *viper.Viper) *config {
	return &config{
		Endpoint:  viper.GetString("huawei.endpoint"),
		Location:  viper.GetString("huawei.location"),
		Bucket:    viper.GetString("huawei.bucket"),
		AccessKey: viper.GetString("huawei.access_key"),
		SecretKey: viper.GetString("huawei.secret_key"),
	}
}
