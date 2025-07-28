package restapis

import (
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/server"
)

type HttpServer struct {
	api_gen.ServerInterface
	App *server.Application
}

func NewHttpServer(app *server.Application) HttpServer {
	return HttpServer{
		App: app,
	}
}
