package request

type InviteToTeamRequest struct {
	UserId int `json:"user_id" binding:"required"`
	TeamId int `json:"team_id" binding:"required"`
}
