package user

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Type int

// Роли
const (
	UserRole Type = iota + 1 // Воскресенье = 0
	OwnerRole
	AdminRole // Вторник = 1, iota автоматически инкрементируется
)

func (t Type) IsValid() bool {
	return t >= UserRole && t <= AdminRole
}

type Role struct {
	Id        int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserId    int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Role      Type      `json:"role" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
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

func GetRolesByUserId(db *gorm.DB, userId int) ([]Role, error) {
	var roles []Role
	result := db.Where("user_id = ?", userId).Find(&roles)

	return roles, result.Error
}
