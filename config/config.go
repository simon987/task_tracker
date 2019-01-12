package config

import (
	"github.com/spf13/viper"
)

var Cfg struct {
	ServerAddr string
	DbConnStr string
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
}
