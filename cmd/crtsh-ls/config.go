package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func configDefaults() {
	viper.SetDefault("crtsh.base_uri", "https://crt.sh/")

	viper.SetDefault("timeout", "60s")

	// viper.SetDefault("command.timeout", "30s")
	viper.RegisterAlias("command.timeout", "general.timeout")
}

func configInit() {
	logrus.SetOutput(os.Stderr)

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Debug("+++configInit()")
	viper.SetConfigName("crtsh")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath(".")

	configDefaults()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Debugf("could not find config file, proceeding with defaults (this is normal)")
	}
}
