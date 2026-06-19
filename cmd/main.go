package main

import (
	"seconda/cmd/config"
	"seconda/cmd/factory"
	"seconda/cmd/service"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()
	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)

	defer dbDecorator.CloseDB()

	factory.BuildAndServe(dbDecorator)
}
