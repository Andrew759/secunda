package config

import (
	config "seconda/pkg"

	"github.com/spf13/viper"
)

type AppConfigurationInterface interface {
	NewAppConfiguration() AppConfiguration
}

type AppConfiguration struct {
	Environment string
	DatabaseConfig
}

func (c AppConfiguration) NewAppConfiguration() AppConfiguration {
	return AppConfiguration{
		DatabaseConfig: PrepareDatabaseConfig(),
	}
}

func PrepareDatabaseConfig() DatabaseConfig {
	dbc := DatabaseConfig{}

	dbc.SetHost(viper.GetString(config.DbHost))
	dbc.SetPort(viper.GetInt(config.DbInternalPort))
	dbc.SetName(viper.GetString(config.DbName))
	dbc.SetUser(viper.GetString(config.DbUser))
	dbc.SetPassword(viper.GetString(config.DbPass))
	dbc.SetTimezone(viper.GetString(config.DbTimezone))

	return dbc
}
