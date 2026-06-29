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

func TestCreateUserWithEmptyPhoneError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Phone' Error:Field validation for 'Phone' failed on the 'required' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithInvalidPhoneError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+7963",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Phone' Error:Field validation for 'Phone' failed on the 'min' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithEmptyNameError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithEmptySurnameError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Surname' Error:Field validation for 'Surname' failed on the 'required' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithEmptyLoginError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "",
		"password": "pass1",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Login' Error:Field validation for 'Login' failed on the 'required' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithEmptyPasswordError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithInvalidPasswordError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "abc",
		"Role": 1
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}

func TestCreateUserWithInvalidRoleError(t *testing.T) {
	tc := base.PrepareTestContainer(t)

	body := []byte(`{
		"phone": "+79634823344",
		"name": "Андрей",
		"surname": "Вельков",
		"login": "avelkov",
		"password": "pass1",
		"Role": 1551
	}`)

	respBytes, statusCode := tc.SendPost(body, "/api/v1/register")

	expectedJSON := `{
		"error": "Invalid input: Key: 'CreateUserRequest.Role' Error:Field validation for 'Role' failed on the 'enum_role' tag"
	}`

	var expectedMap map[string]any
	json.Unmarshal([]byte(expectedJSON), &expectedMap)

	var actualMap map[string]any
	json.Unmarshal(respBytes, &actualMap)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, expectedMap, actualMap)
}
