package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg struct {
	ServerAddr       string
	DbConnStr        string
	WebHookSecret    []byte
	WebHookHash      string
	WebHookSigHeader string
	LogLevel         logrus.Level
	DbLogLevels      []logrus.Level
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
	Cfg.LogLevel, _ = logrus.ParseLevel(viper.GetString("log.level"))
	for _, level := range viper.GetStringSlice("database.log_levels") {
		newLevel, _ := logrus.ParseLevel(level)
		Cfg.DbLogLevels = append(Cfg.DbLogLevels, newLevel)
	}
}
