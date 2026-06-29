package user

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone     string    `json:"phone" gorm:"type:varchar(30)" binding:"required,min=7,max=15"`
	Name      string    `json:"name" gorm:"type:varchar(256);not null" binding:"required,min=1,max=256"`
	Surname   string    `json:"surname" gorm:"type:varchar(256);not null" binding:"required,min=1,max=256"`
	Login     string    `json:"login" gorm:"type:varchar(256);unique; not null" binding:"required,min=1,max=256"`
	Password  string    `json:"-" gorm:"type:varchar(256)" binding:"required,min=5,max=256"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}

var WithPhoneAlreadyExistsErr = errors.New("user with phone already exist")

var WithLoginAlreadyExistsErr = errors.New("user with login already exist")

var NotFoundErr = errors.New("user not found")

func (u User) TableName() string {
	return "users"
}

func passwordHash(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func CreateUser(ctx context.Context, db *gorm.DB, u *User) error {
	var existing User

	err := db.WithContext(ctx).Where("login = ? OR phone = ?", u.Login, u.Phone).First(&existing).Error

	if err == nil {
		if existing.Login == u.Login {
			return WithLoginAlreadyExistsErr
		}
		if existing.Phone == u.Phone {
			return WithPhoneAlreadyExistsErr
		}
		return err
	}

	if errors.Is(err, context.Canceled) {
		return err
	}

	password, err := passwordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = password

	return db.WithContext(ctx).Create(u).Error
}

func GetUserByLoginAndPass(ctx context.Context, db *gorm.DB, login, password string) (User, error) {
	var user User

	result := db.WithContext(ctx).Where("login = ?", login).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, NotFoundErr
		}
		return user, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}
