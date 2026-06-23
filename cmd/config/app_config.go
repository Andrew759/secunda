package config

import (
	"seconda/pkg/config"

	"github.com/spf13/viper"
)

type AppConfigurationInterface interface {
	NewAppConfiguration() AppConfiguration
}

type AppConfiguration struct {
	DatabaseConfig
	RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       int
}

func (c AppConfiguration) NewAppConfiguration() AppConfiguration {
	return AppConfiguration{
		DatabaseConfig: PrepareDatabaseConfig(),
		RedisConfig:    PrepareRedisConfig(),
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

func PrepareRedisConfig() RedisConfig {
	rc := RedisConfig{}

	rc.Host = viper.GetString(config.RedisHost)
	rc.Port = viper.GetInt(config.RedisInternalPort)
	rc.User = viper.GetString(config.RedisUser)
	rc.Password = viper.GetString(config.RedisPassword)
	rc.Db = viper.GetInt(config.RedisDB)

	return rc
}
