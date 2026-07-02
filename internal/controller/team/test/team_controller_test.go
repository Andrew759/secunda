package test

import (
	"encoding/json"
	"net/http"
	base "seconda/internal/controller/test"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateTeamSuccess(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	registerBody := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	_, cookies, _ := tc.SendPostWithAuthToken(registerBody, "/api/v1/register", "")

	var accessToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessToken = cookie.Value
			break
		}
	}

	teamBody := []byte(`{
		"name": "команда"
	}`)

	respBytes, _, statusCode := tc.SendPostWithAuthToken(teamBody, "/api/v1/teams", accessToken)

	var actual struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		CreatedBy int    `json:"created_by"`
	}

	json.Unmarshal(respBytes, &actual)

	assert.Equal(t, http.StatusOK, statusCode) // Обычно для создания используют 201 Created, проверьте ваш API
	assert.Equal(t, 1, actual.Id)
	assert.Equal(t, "команда", actual.Name)
	assert.Equal(t, 1, actual.CreatedBy)
}

func TestGetTeamsSuccess(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	registerBody := []byte(`{"phone": "+79634823344", "name": "Андрей", "surname": "Вельков", "login": "avelkov", "password": "pass1", "Role": 1}`)
	_, cookies, _ := tc.SendPostWithAuthToken(registerBody, "/api/v1/register", "")
	var accessToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessToken = cookie.Value
			break
		}
	}

	teamBody := []byte(`{"name": "Моя Команда"}`)
	tc.SendPostWithAuthToken(teamBody, "/api/v1/teams", accessToken)

	inviteBody := []byte(`{"user_id": 1}`)
	tc.SendPostWithAuthToken(inviteBody, "/api/v1/teams/1/invite", accessToken)

	respBytes, _, statusCode := tc.SendGetWithAuthToken("/api/v1/teams", accessToken)

	var actual []struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		CreatedBy int    `json:"created_by"`
	}
	json.Unmarshal(respBytes, &actual)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, "Моя Команда", actual[0].Name)
}

func TestGetTeamsUnauthorized(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	_, _, statusCode := tc.SendGetWithAuthToken("/api/v1/teams", "")

	assert.Equal(t, http.StatusUnauthorized, statusCode)
}
