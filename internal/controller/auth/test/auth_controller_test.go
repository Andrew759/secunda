package test

import (
	"encoding/json"
	"net/http"
	base "seconda/internal/controller/test"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateUserSuccess(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"id": 1,
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	//Не сравниваемые поля
	delete(actualMap, "created_at")
	delete(actualMap, "updated_at")

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}
