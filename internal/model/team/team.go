package team

import (
	"errors"
	"os/user"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	Id            int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null;"`
	CreatedBy     int       `json:"created_by" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser user.User `json:"created_by_user" gorm:"foreignKey:createdBy;references:Id"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var NotFoundErr = errors.New("team not found")

func (t Team) TableName() string {
	return "teams"
}

func CreateTeam(db *gorm.DB, t *Team) error {
	return db.Create(t).Error
}

func GetTeamById(db *gorm.DB, id int) (Team, error) {
	var team Team
	err := db.Where("id = ?", id).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Team{}, NotFoundErr
	}

	return team, err
}
