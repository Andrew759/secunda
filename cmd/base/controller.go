package base

import (
	"seconda/cmd/service"

	"github.com/gin-gonic/gin"
)

type DIContainer struct {
	DBDecorator    *service.DBDecorator
	RedisDecorator *service.RedisDecorator
}

type Controller struct {
	E  *gin.Engine
	DI *DIContainer
}

type RequestHandler interface {
	HandleRequest()
}
