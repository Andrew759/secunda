package main

import (
	"log/slog"
	"seconda/cmd/config"
	"seconda/cmd/factory"
	"seconda/cmd/service"
	"seconda/internal/model/task"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	//TODO: не коммитить!
	dbDecorator.GDB().AutoMigrate(user.User{}, user.Role{}, team.Team{}, team.Member{}, task.Task{}, task.Comment{}, task.History{})

	redisDecorator := service.InitRedis(appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	err := factory.BuildAndServe(dbDecorator, redisDecorator)
	if err != nil {
		slog.Error(err.Error())
	}
}
