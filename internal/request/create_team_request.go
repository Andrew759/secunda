package request

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=2,max=256"`
}
