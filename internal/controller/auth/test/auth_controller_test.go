package test

import (
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

	_, statusCode := tc.SendPost(body, "/api/v1/register")

	assert.Equal(t, http.StatusCreated, statusCode)

	//expectedJSON := `{
	//	"phone": "+79634823344",
	//	"name": "Андрей",
	//	"surname": "Вельков",
	//	"login": "avelkov",
	//	"password": "pass1",
	//	"Role": 1
	//}`
	//
	//var expectedMap map[string]any
	//if err := json.Unmarshal([]byte(expectedJSON), &expectedMap); err != nil {
	//	t.Fatalf("Ошибка парсинга ожидаемого JSON: %v", err)
	//}
	//
	//var actualMap map[string]any
	//if err := json.Unmarshal(respBytes, &actualMap); err != nil {
	//	t.Fatalf("Ошибка парсинга ответа сервера: %v. Сырой ответ: %s", err, string(respBytes))
	//}
	//
	//assert.Equal(t, expectedMap, actualMap)
}
