package user

import (
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
	Password  string    `json:"password" gorm:"type:varchar(256)" binding:"required,min=5,max=256"`
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

func CreateUser(db *gorm.DB, u *User) error {
	var existing User
	err := db.Where("login = ? OR phone = ?", u.Login, u.Phone).First(&existing).Error

	if err == nil {
		if existing.Login == u.Login {
			return WithLoginAlreadyExistsErr
		}
		if existing.Phone == u.Phone {
			return WithPhoneAlreadyExistsErr
		}
		return err
	}

	password, err := passwordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = password

	return db.Create(u).Error
}

// GetUserById TODO: удалить, если не будет использоваться
func GetUserById(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, NotFoundErr
	}

	return user, result.Error
}

func GetUserByLoginAndPass(db *gorm.DB, login, password string) (User, error) {
	var user User

	result := db.Where("login = ?", login).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, NotFoundErr
	}

	return user, result.Error
}
