package team

import (
	"context"
	"errors"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	Id            int        `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	Name          string     `json:"name" gorm:"type:varchar(256);not null" binding:"required,min=2,max=256"`
	CreatedBy     int        `json:"created_by" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" binding:"required"`
	CreatedByUser *user.User `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:timestamp;not null"`
}

var NotFoundErr = errors.New("team not found")

var AlreadyExistErr = errors.New("team already exist")

func (t Team) TableName() string {
	return "teams"
}

func CreateTeam(ctx context.Context, db *gorm.DB, t *Team) error {
	var existing Team

	err := db.WithContext(ctx).Where("name = ? AND created_by = ?", t.Name, t.CreatedBy).First(&existing).Error

	if err == nil {
		if existing.Name == t.Name && existing.CreatedBy == t.CreatedBy {
			return AlreadyExistErr
		}
		return err
	}

	return db.WithContext(ctx).Create(t).Error
}

func GetTeamById(ctx context.Context, db *gorm.DB, id int) (Team, error) {
	var team Team
	err := db.WithContext(ctx).Where("id = ?", id).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Team{}, NotFoundErr
	}

	return team, err
}

func GetTeamsWhereUserIsMember(ctx context.Context, db *gorm.DB, userId int) ([]Team, error) {
	var teams []Team

	err := db.WithContext(ctx).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userId).
		Find(&teams).Error

	return teams, err
}
