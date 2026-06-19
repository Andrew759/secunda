package factory

import (
	"net/http"
	"seconda/cmd/base"
	"seconda/cmd/service"
)

func BuildAndServe(dbDecorator service.DBDecorator) {
	mux := BuildServer(dbDecorator)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator) *http.ServeMux {
	mux := http.NewServeMux()

	diContainer := base.DIContainer{
		DBDecorator: dbDecorator,
	}
	initCommandWorkService(mux, diContainer)

	return mux
}

func initCommandWorkService(mux *http.ServeMux, diContainer base.DIContainer) {

}
