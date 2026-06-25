package team

import (
	"errors"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type Member struct {
	UserId    int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      user.User `json:"user" gorm:"foreignKey:UserId;references:Id"`
	TeamId    int       `json:"team_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Team      Team      `json:"team" gorm:"foreignKey:TeamId;references:Id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var MemberNotFoundErr = errors.New("team member not found")

var WithUserIdAlreadyExistErr = errors.New("user with id already exist")

var WithTeamIdAlreadyExistErr = errors.New("user with team id already exist")

func (tm Member) TableName() string {
	return "team_members"
}

func CreateMember(db *gorm.DB, m *Member) error {
	var existing Member
	err := db.Where("user_id = ? AND team_id = ?", m.UserId, m.TeamId).First(&existing).Error

	if err == nil {
		if existing.UserId == m.UserId {
			return WithUserIdAlreadyExistErr
		}
		if existing.TeamId == m.TeamId {
			return WithTeamIdAlreadyExistErr
		}
		return err
	}
	return db.Create(m).Error
}

func GetMemberById(db *gorm.DB, id int) (Member, error) {
	var member Member
	err := db.Where("id = ?", id).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Member{}, MemberNotFoundErr
	}

	return member, err
}
