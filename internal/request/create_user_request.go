package request

import "seconda/internal/model/user"

type CreateUserRequest struct {
	Phone    string    `json:"phone" binding:"required,min=7,max=15"`
	Name     string    `json:"name" binding:"required,min=1,max=256"`
	Surname  string    `json:"surname" binding:"required,min=1,max=256"`
	Login    string    `json:"login" binding:"required,min=1,max=256"`
	Password string    `json:"password" binding:"required,min=5,max=256"`
	Role     user.Type `json:"role" binding:"required,enum_role"`
}
