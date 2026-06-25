package task

import (
	"errors"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	Id            int       `json:"id" gorm:"type:int;not null;primaryKey;autoIncrement"`
	AssigneeId    int       `json:"assignee_id" gorm:"type:int;not null;" binding:"required"`
	AssigneeUser  user.User `json:"assignee_user" gorm:"foreignKey:AssigneeId;references:Id"`
	TeamId        int       `json:"team_id" gorm:"type:int;not null;" binding:"required"`
	Team          team.Team `json:"team" gorm:"foreignKey:TeamId;references:Id"`
	CreatedBy     int       `json:"created_by" gorm:"type:int;not null;" binding:"required"`
	CreatedByUser user.User `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var NotFoundErr = errors.New("task not found")

func (t Task) TableName() string {
	return "tasks"
}

func CreateTask(db *gorm.DB, t *Task) error {
	return db.Create(t).Error
}

func GetTaskById(db *gorm.DB, id int) (Task, error) {
	var t Task
	err := db.Where("id = ?", id).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Task{}, NotFoundErr
	}

	return t, nil
}
