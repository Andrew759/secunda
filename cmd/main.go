package main

import (
	"log/slog"
	"seconda/cmd/config"
	"seconda/cmd/factory"
	"seconda/cmd/service"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	redisDecorator := service.InitRedis(&appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	err := factory.BuildAndServe(dbDecorator, redisDecorator)
	if err != nil {
		slog.Error(err.Error())
	}
}
