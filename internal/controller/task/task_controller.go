package task

import (
	"seconda/cmd/base"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Controller base.Controller
}

func (tc *TaskController) HandleRequest() {
	e := tc.Controller.E

	group := e.Group("/api/v1")
	group.POST("/tasks", tc.CreateTask)
	group.GET("/tasks", tc.Tasks)
	group.PUT("/tasks/:id", tc.UpdateTask)
	group.GET("/tasks/:id/history", tc.TaskHistory)
}

func (tc *TaskController) CreateTask(c *gin.Context) {

}

func (tc *TaskController) Tasks(c *gin.Context) {

}

func (tc *TaskController) UpdateTask(c *gin.Context) {

}

func (tc *TaskController) TaskHistory(c *gin.Context) {

}
