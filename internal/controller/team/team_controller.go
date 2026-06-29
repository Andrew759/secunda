package team

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"seconda/cmd/base"
	"seconda/internal/middleware"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"seconda/internal/request"
	"strconv"
	"time"

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
	group.GET("/teams/report", tc.GetTeamsReportCSV)
}

func (tc *TeamController) CreateTeams(c *gin.Context) {
	ctx := c.Request.Context()

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
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))
	t.CreatedBy = userIdI

	if err := team.CreateTeam(ctx, tc.Controller.DI.DBDecorator.GDB(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc *TeamController) Teams(c *gin.Context) {
	ctx := c.Request.Context()

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}
	userIdI, _ := strconv.Atoi(userId.(string))

	teams, err := team.GetTeamsWhereUserIsMember(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (tc *TeamController) Invite(c *gin.Context) {
	ctx := c.Request.Context()

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

	teamIdStr := c.Param("id")
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := team.GetTeamById(ctx, tc.Controller.DI.DBDecorator.GDB(), teamId)
	if err != nil && errors.Is(err, team.NotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, err := user.GetRolesByUserId(ctx, tc.Controller.DI.DBDecorator.GDB(), userIdI)
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
	m.TeamId = teamId

	err = team.CreateMember(ctx, tc.Controller.DI.DBDecorator.GDB(), &m)
	if err != nil && errors.Is(err, team.MemberAlreadyExist) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

func (tc *TeamController) GetTeamsReportCSV(c *gin.Context) {
	ctx := c.Request.Context()

	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
		return
	}

	reports, err := team.GetTeamsActivityReport(ctx, tc.Controller.DI.DBDecorator.GDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Team ID", "Team Name", "Members Count", "Done Tasks (Last 7 Days)"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers: " + err.Error()})
		return
	}

	for _, report := range reports {
		row := []string{
			strconv.Itoa(report.TeamId),
			report.TeamName,
			strconv.Itoa(report.MembersCount),
			strconv.Itoa(report.DoneTasksCount),
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data: " + err.Error()})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flush CSV writer: " + err.Error()})
		return
	}

	fileName := fmt.Sprintf("teams_activity_report_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}
