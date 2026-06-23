package factory

import (
	"seconda/cmd/base"
	"seconda/cmd/service"
	"seconda/internal/controller/auth"

	"github.com/gin-gonic/gin"
)

func BuildAndServe(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) error {
	e := gin.Default()

	initAuthService(e, &base.DIContainer{DBDecorator: dbDecorator, RedisDecorator: redisDecorator})

	err := e.Run()
	if err != nil {
		return err
	}

	return nil
}

func initAuthService(e *gin.Engine, aDIC *base.DIContainer) {
	authService := auth.AuthController{
		Controller: base.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	authService.HandleRequest()
}
