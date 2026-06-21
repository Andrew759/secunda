package user

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Type int

// TODO: убедиться что роли будут такими
const (
	UserRole  Type = iota // Воскресенье = 0
	OwnerRole             // Вторник = 1, iota автоматически инкрементируется
)

type Role struct {
	//TODO: убедиться, что правило unique корректно
	UserId    int       `json:"user_id" gorm:"type:int;unique;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Role      Type      `json:"role" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"type:timestamp;not null"`
}

var RoleNotFound = errors.New("role not found")

func (r Role) TableName() string {
	return "roles"
}

func CreateRole(db *gorm.DB, r *Role) error {
	return db.Create(r).Error
}

func GetRoleById(db *gorm.DB, id int) (Role, error) {
	var role Role
	result := db.First(&role, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Role{}, RoleNotFound
	}

	return role, result.Error
}
