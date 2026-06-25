package team

import (
	"errors"
	"net/http"
	"seconda/cmd/base"
	"seconda/internal/middleware"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"seconda/internal/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	Controller base.Controller
}

func (tc *TeamController) HandleRequest() {
	e := tc.Controller.E

	group := e.Group("/api/v1")
	group.Use(middleware.AuthMiddleware())

	group.POST("/teams", tc.CreateTeams)
	group.GET("/teams", tc.Teams)
	group.POST("/teams/:id/invite", tc.Invite)
}

func (tc *TeamController) CreateTeams(c *gin.Context) {
	var ctr request.CreateTeamRequest
	if err := c.ShouldBindJSON(&ctr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var t team.Team
	t.Name = ctr.Name
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
	}
	userIdI, _ := strconv.Atoi(userId.(string))
	t.CreatedBy = userIdI

	if err := team.CreateTeam(tc.Controller.DI.DBDecorator.GDB(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc *TeamController) Teams(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	teams, err := team.GetTeamsWhereUserIsMember(tc.Controller.DI.DBDecorator.GDB(), userIdI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (tc *TeamController) Invite(c *gin.Context) {
	var ittr request.InviteToTeamRequest
	if err := c.ShouldBindJSON(&ittr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := team.GetTeamById(tc.Controller.DI.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, team.NotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, err := user.GetRolesByUserId(tc.Controller.DI.DBDecorator.GDB(), userIdI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	hasAccessByRole := false
	for _, role := range roles {
		if role.Role == user.OwnerRole || role.Role == user.AdminRole {
			hasAccessByRole = true
		}
	}
	if t.CreatedBy != userIdI && !hasAccessByRole {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "you cannot invite a user if you are not the owner or admin"})
		return
	}

	var m team.Member
	m.UserId = ittr.UserId
	m.TeamId = ittr.TeamId

	err = team.CreateMember(tc.Controller.DI.DBDecorator.GDB(), &m)
	if err != nil && errors.Is(err, team.WithUserIdAlreadyExistErr) || errors.Is(err, user.WithLoginAlreadyExistsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}
