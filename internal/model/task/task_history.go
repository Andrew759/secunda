package task

import (
	"errors"
	"os/user"
	"time"

	"gorm.io/gorm"
)

type History struct {
	Id            int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	TaskId        int       `json:"task_id" gorm:"type:int;not null;"`
	Task          Task      `json:"task" gorm:"foreignKey:taskId;references:Id"`
	ChangedBy     int       `json:"changed_by" gorm:"type:int;not null;"`
	ChangedByUser user.User `json:"changed_by_user" gorm:"foreignKey:changedBy;references:Id"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"type:timestamp;not null"`
}

var HistoryNotFoundErr = errors.New("task history not found")

func (h History) TableName() string {
	return "task_history"
}

func CreateHistory(db *gorm.DB, h *History) error {
	return db.Create(h).Error
}

func GetHistoryById(db *gorm.DB, id int) (History, error) {
	var history History
	err := db.Where("id = ?", id).First(&history).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return History{}, HistoryNotFoundErr
	}

	return history, err
}
