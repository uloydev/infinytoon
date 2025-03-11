package main

import (
	"github.com/spf13/viper"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
)

type Config struct {
	AppName    string                      `mapstructure:"app_name"`
	AppVersion string                      `mapstructure:"app_version"`
	AppEnv     string                      `mapstructure:"app_env"`
	Server     quictunnel.QuicServerConfig `mapstructure:"server"`
	Logger     logger.LoggerConfig         `mapstructure:"logger"`
}

func InitConfig(appCtx *appctx.AppContext) *Config {
	cfg := Config{}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	appCtx.Set(appctx.ConfigKey, &cfg)
	appCtx.Set(appctx.AppNameKey, cfg.AppName)
	appCtx.Set(appctx.EnvironmentKey, cfg.AppEnv)

	return &cfg
}
