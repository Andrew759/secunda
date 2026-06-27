package task

import (
	"context"
	"errors"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type History struct {
	Id            int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	TaskId        int       `json:"task_id" gorm:"type:int;not null;"`
	Task          Task      `json:"task" gorm:"foreignKey:taskId;references:Id"`
	ChangedBy     int       `json:"changed_by" gorm:"type:int;not null;"`
	ChangedByUser user.User `json:"changed_by_user" gorm:"foreignKey:ChangedBy;references:Id"`
	TeamId        int       `json:"team_id" gorm:"type:int;not null;"`
	Team          team.Team `json:"team" gorm:"foreignKey:TeamId;references:Id"`
	CreatedBy     int       `json:"created_by" gorm:"type:int;not null;"`
	CreatedByUser user.User `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	Name          string    `json:"name" gorm:"type:varchar(256);not null;"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var HistoryNotFoundErr = errors.New("task history not found")

func (h History) TableName() string {
	return "task_history"
}

func CreateHistory(ctx context.Context, db *gorm.DB, h *History) error {
	return db.WithContext(ctx).Create(h).Error
}

func GetHistoryListByTaskId(ctx context.Context, db *gorm.DB, taskId int) ([]History, error) {
	var history []History

	err := db.WithContext(ctx).
		Where("task_id = ?", taskId).
		Order("id ASC").
		Find(&history).Error

	if err != nil {
		return nil, err
	}

	return history, nil
}
