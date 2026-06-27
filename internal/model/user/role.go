package user

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id        int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserId    int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Role      Type      `json:"role" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

type Type int

const (
	UserRole Type = iota + 1
	OwnerRole
	AdminRole
)

func (t Type) IsValid() bool {
	return t >= UserRole && t <= AdminRole
}

var RoleNotFound = errors.New("role not found")

func (r Role) TableName() string {
	return "roles"
}

func CreateRole(ctx context.Context, db *gorm.DB, r *Role) error {
	return db.WithContext(ctx).Create(r).Error
}

func GetRolesByUserId(ctx context.Context, db *gorm.DB, userId int) ([]Role, error) {
	var roles []Role
	result := db.WithContext(ctx).Where("user_id = ?", userId).Find(&roles)

	return roles, result.Error
}
