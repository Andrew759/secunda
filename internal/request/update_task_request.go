package request

import (
	"seconda/internal/enum"
)

type UpdateTaskRequest struct {
	AssigneeId int         `json:"assignee_id" binding:"required"`
	TeamId     int         `json:"team_id" binding:"required"`
	CreatedBy  int         `json:"created_by" binding:"required"`
	Status     enum.Status `json:"status" binding:"required,task_status"`
	Name       string      `json:"name" binding:"required,min=2,max=256"`
	Comment    *string     `json:"comment" binding:"omitempty,min=2,max=1000"`
}
