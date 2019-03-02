package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var Cfg struct {
	ServerAddr                 string
	DbConnStr                  string
	WebHookSecret              []byte
	WebHookHash                string
	WebHookSigHeader           string
	LogLevel                   logrus.Level
	DbLogLevels                []logrus.Level
	SessionCookieName          string
	SessionCookieExpiration    time.Duration
	MonitoringInterval         time.Duration
	ResetTimedOutTasksInterval time.Duration
	MonitoringHistory          time.Duration
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
	Cfg.SessionCookieName = viper.GetString("session.cookie_name")
	Cfg.SessionCookieExpiration, err = time.ParseDuration(viper.GetString("session.expiration"))
	Cfg.MonitoringInterval, err = time.ParseDuration(viper.GetString("monitoring.snapshot_interval"))
	handleErr(err)
	Cfg.ResetTimedOutTasksInterval, err = time.ParseDuration(viper.GetString("maintenance.reset_timed_out_tasks_interval"))
	handleErr(err)
	Cfg.MonitoringHistory, err = time.ParseDuration(viper.GetString("monitoring.history_length"))
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
