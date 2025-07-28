package restapis

import (
	"github.com/slilp/go-wallet/internal/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/server"
)

type HttpServer struct {
	api_gen.ServerInterface
	app *server.Application
}

func NewHttpServer(app *server.Application) HttpServer {
	return HttpServer{
		app: app,
	}
}
