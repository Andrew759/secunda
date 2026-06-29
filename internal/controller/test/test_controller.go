package base

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"seconda/cmd/config"
	"seconda/cmd/factory"
	"seconda/cmd/service"
	"seconda/internal/model/task"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestContainer struct {
	service.DBDecorator
	HTTPServer *httptest.Server
	HTTPClient *http.Client
}

func PrepareTestContainer(t *testing.T) TestContainer {
	t.Helper()

	gin.SetMode(gin.TestMode)

	factory.InitViper()

	appCfg := config.AppConfiguration{}.NewAppConfiguration()
	dbd := service.InitORM(&appCfg.DatabaseConfig)
	rd := service.InitRedis(&appCfg.RedisConfig)

	router := factory.BuildServer(dbd, rd)
	server := httptest.NewServer(router)

	CreateTables(dbd)

	t.Cleanup(func() {
		DropTables(dbd)
		dbd.CloseDB()
		server.Close()
	})

	return TestContainer{
		DBDecorator: *dbd,
		HTTPServer:  server,
		HTTPClient:  InitHttpClient(),
	}
}

func GetTableName(db *gorm.DB, model any) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func InitHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
			ResponseHeaderTimeout: time.Second,
		},
	}
}

func CreateTables(db *service.DBDecorator) {
	if err := db.GDB().AutoMigrate(
		user.User{},
		user.Role{},
		team.Team{},
		team.Member{},
		task.Task{},
		task.Comment{},
		task.History{},
	); err != nil {
		panic("failed to create tables: " + err.Error())
	}
}

func DropTables(db *service.DBDecorator) {
	models := []any{
		user.User{}, user.Role{}, team.Team{},
		team.Member{}, task.Task{}, task.Comment{}, task.History{},
	}

	tableNames := make([]string, 0, len(models))
	for _, model := range models {
		tName, err := GetTableName(db.GDB(), model)
		if err != nil {
			panic("failed to get table name: " + err.Error())
		}
		tableNames = append(tableNames, tName)
	}

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", strings.Join(tableNames, ", "))
	if _, err := db.NativeDB().Exec(query); err != nil {
		panic("failed to drop tables: " + err.Error())
	}
}

func (tc TestContainer) SendPost(body []byte, uri string) ([]byte, int) {
	resp, err := tc.HTTPClient.Post(
		tc.HTTPServer.URL+uri,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, 0
	}
	defer resp.Body.Close()

	respBodyBytes, _ := io.ReadAll(resp.Body)

	return respBodyBytes, resp.StatusCode
}
