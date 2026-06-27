package task

import (
	"errors"
	"net/http"
	"seconda/cmd/base"
	"seconda/internal/middleware"
	"seconda/internal/model/task"
	"seconda/internal/model/team"
	"seconda/internal/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Controller base.Controller
}

func (tc *TaskController) HandleRequest() {
	e := tc.Controller.E

	group := e.Group("/api/v1")
	group.Use(middleware.AuthMiddleware())

	group.POST("/tasks", tc.CreateTask)
	group.GET("/tasks", tc.Tasks)
	group.PUT("/tasks/:id", tc.UpdateTask)
	group.GET("/tasks/:id/history", tc.TaskHistory)
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var ctr request.CreateTaskRequest
	if err := c.ShouldBindJSON(&ctr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var t task.Task
	t.AssigneeId = ctr.AssigneeId
	t.TeamId = ctr.TeamId
	t.Name = ctr.Name
	t.CreatedBy = ctr.CreatedBy
	t.Status = ctr.Status

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	_, err := team.GetMemberByUserIdAndTeamId(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI, ctr.TeamId)
	if err != nil && errors.Is(err, team.MemberNotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "current user is not a team member"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = task.CreateTask(ctx, tc.Controller.DI.DBDecorator.GDB(), &t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc *TaskController) Tasks(c *gin.Context) {
	ctx := c.Request.Context()

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	var filter task.Filter

	if teamIdStr := c.Query("team_id"); teamIdStr != "" {
		teamId, err := strconv.Atoi(teamIdStr)
		if err != nil || teamId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team_id parameter"})
			return
		}
		filter.TeamId = teamId
	}

	if assigneeIdStr := c.Query("assignee_id"); assigneeIdStr != "" {
		assigneeId, err := strconv.Atoi(assigneeIdStr)
		if err != nil || assigneeId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id parameter"})
			return
		}
		filter.AssigneeId = assigneeId
	}

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil && statusStr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status parameter"})
		return
	}

	if status != 0 {
		switch status {
		case int(task.Draft):
			filter.Status = int(task.Draft)
		case int(task.Todo):
			filter.Status = int(task.Todo)
		case int(task.InProgress):
			filter.Status = int(task.InProgress)
		case int(task.Done):
			filter.Status = int(task.Done)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status parameter"})
			return
		}
	}

	if filter.TeamId > 0 {
		_, err = team.GetMemberByUserIdAndTeamId(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI, filter.TeamId)
		if err != nil && errors.Is(err, team.MemberNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": "current user is not a team member"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	tasks, err := task.GetTasksByFilter(ctx, tc.Controller.DI.DBDecorator.GDB(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks: " + err.Error()})
		return
	}

	if tasks == nil {
		tasks = []task.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var utr request.UpdateTaskRequest
	if err := c.ShouldBindJSON(&utr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var nt task.Task
	nt.AssigneeId = utr.AssigneeId
	nt.Status = utr.Status
	nt.Name = utr.Name
	nt.TeamId = utr.TeamId
	nt.CreatedBy = utr.CreatedBy

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	_, err := team.GetMemberByUserIdAndTeamId(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI, utr.TeamId)
	if err != nil && errors.Is(err, team.MemberNotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "current user is not a team member"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	ot, err := task.UpdateTaskById(ctx, tc.Controller.DI.DBDecorator.GDB(), &nt, taskId)
	if err != nil {
		if errors.Is(err, task.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var h task.History
	h.TaskId = ot.Id
	h.ChangedBy = userIdI
	h.TeamId = ot.TeamId
	h.CreatedBy = ot.CreatedBy
	h.Name = ot.Name

	err = task.CreateHistory(ctx, tc.Controller.DI.DBDecorator.GDB(), &h)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, nt)
}

func (tc *TaskController) TaskHistory(c *gin.Context) {
	ctx := c.Request.Context()

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	t, err := task.GetTaskById(ctx, tc.Controller.DI.DBDecorator.GDB(), taskId)
	if err != nil {
		if errors.Is(err, task.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = team.GetMemberByUserIdAndTeamId(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI, t.TeamId)
	if err != nil {
		if errors.Is(err, team.MemberNotFoundErr) {
			c.JSON(http.StatusForbidden, gin.H{"error": "current user is not a member of the task's team"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	history, err := task.GetHistoryListByTaskId(ctx, tc.Controller.DI.DBDecorator.GDB(), taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task history: " + err.Error()})
		return
	}

	if history == nil {
		history = []task.History{}
	}

	c.JSON(http.StatusOK, history)
}
