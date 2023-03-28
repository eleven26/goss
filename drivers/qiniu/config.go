package qiniu

import "github.com/spf13/viper"

type config struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Zone      string
	Domain    string
	Private   bool
}

func getConfig(viper *viper.Viper) *config {
	return &config{
		Bucket:    viper.GetString("qiniu.bucket"),
		AccessKey: viper.GetString("qiniu.access_key"),
		SecretKey: viper.GetString("qiniu.secret_key"),
		Zone:      viper.GetString("qiniu.zone"),
		Domain:    viper.GetString("qiniu.domain"),
		Private:   viper.GetBool("qiniu.private"),
	}
}
