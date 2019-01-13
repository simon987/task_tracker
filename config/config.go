package config

import (
	"github.com/spf13/viper"
)

var Cfg struct {
	ServerAddr       string
	DbConnStr        string
	WebHookSecret    []byte
	WebHookHash      string
	WebHookSigHeader string
}

func SetupConfig() {

	viper.AddConfigPath(".")
	viper.SetConfigName("")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Cfg.ServerAddr = viper.GetString("server.address")
	Cfg.DbConnStr = viper.GetString("database.conn_str")
	Cfg.WebHookSecret = []byte(viper.GetString("git.webhook_secret"))
	Cfg.WebHookHash = viper.GetString("git.webhook_hash")
	Cfg.WebHookSigHeader = viper.GetString("git.webhook_sig_header")
}
