package factory

import (
	"seconda/cmd/base"
	"seconda/cmd/service"
	"seconda/internal/controller/auth"
	taskService "seconda/internal/controller/task"
	"seconda/internal/controller/team"
	"seconda/internal/enum"
	"seconda/internal/model/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func BuildAndServe(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) error {
	e := BuildServer(dbDecorator, redisDecorator)

	return e.Run()
}

func BuildServer(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) *gin.Engine {
	e := gin.Default()

	RegisterRoleValidator()
	RegisterStatusValidator()

	diContainer := &base.DIContainer{
		DBDecorator:    dbDecorator,
		RedisDecorator: redisDecorator,
	}

	initAuthService(e, diContainer)
	initTeamsService(e, diContainer)
	InitTasksService(e, diContainer)

	return e
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

func initTeamsService(e *gin.Engine, aDIC *base.DIContainer) {
	teamController := team.TeamController{
		Controller: base.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	teamController.HandleRequest()
}

func InitTasksService(e *gin.Engine, aDIC *base.DIContainer) {
	taskController := taskService.TaskController{
		Controller: base.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	taskController.HandleRequest()
}

// RegisterRoleValidator регистрирует валидатор глобально для всего движка Gin.
func RegisterRoleValidator() bool {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("enum_role", func(fl validator.FieldLevel) bool {
			if val, ok := fl.Field().Interface().(user.Type); ok {
				return val.IsValid()
			}
			return false
		})
		return err == nil
	}
	return false
}

func RegisterStatusValidator() bool {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("task_status", func(fl validator.FieldLevel) bool {
			if val, ok := fl.Field().Interface().(enum.Status); ok {
				return val.IsValid()
			}
			return false
		})
		return err == nil
	}
	return false
}
