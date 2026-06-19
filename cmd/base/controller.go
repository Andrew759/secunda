package base

import (
	"net/http"
	"seconda/cmd/service"
)

type DIContainer struct {
	DBDecorator service.DBDecorator
}

type Controller struct {
	ServeMux     *http.ServeMux
	Dependencies DIContainer
}

type RequestHandler interface {
	HandleRequest()
}
