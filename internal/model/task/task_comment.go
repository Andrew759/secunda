package task

import (
	"errors"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id        int       `json:"id" gorm:"type:int;not null;primaryKey;autoIncrement"`
	TaskId    int       `json:"task_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Task      Task      `json:"task" gorm:"foreignKey:taskId;references:Id"`
	UserId    int       `json:"user_id" gorm:"type:int;unique;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      user.User `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Comment   string    `json:"comment" gorm:"type:text;not null;"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var CommentNotFoundErr = errors.New("task comment not found")

func (c Comment) TableName() string {
	return "task_comments"
}

func CreateComment(db *gorm.DB, c *Comment) error {
	return db.Create(c).Error
}

func GetCommentById(db *gorm.DB, id int) (Comment, error) {
	var c Comment
	err := db.Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Comment{}, CommentNotFoundErr
	}

	return c, nil
}
