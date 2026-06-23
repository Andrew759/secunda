package request

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required,min=2,max=256"`
	Password string `json:"password" binding:"required,min=5,max=256"`
}
