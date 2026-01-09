package server

import (
	"context"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/service"
	"github.com/kibirisu/borg/internal/worker"
	"github.com/kibirisu/borg/web"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct {
	assets  fs.FS
	conf    *config.Config
	service *service.Container
	worker  worker.Worker
}

func New(ctx context.Context, conf *config.Config) *http.Server {
	assets := web.GetAssets()
	service := service.NewContainer(ctx, conf)
	worker := worker.New(ctx)
	server := &Server{
		assets,
		conf,
		service,
		worker,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(preAuthMiddleware)
	r.Group(server.federationRoutes())

	h := api.HandlerWithOptions(
		server,
		api.ChiServerOptions{
			BaseRouter:  r,
			Middlewares: []api.MiddlewareFunc{server.createAuthMiddleware()},
		},
	)
	r.Group(server.staticRoutes())

	s := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              "0.0.0.0:" + conf.ListenPort,
	}
	return s
}
