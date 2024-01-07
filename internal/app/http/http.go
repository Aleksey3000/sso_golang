package http

import (
	"SSO/internal/config"
	"SSO/internal/http/apps"
	"fmt"
)

type App struct {
	server *apps.HttpServer
}

func NewHttpApp(appsServer apps.Apps, cnf *config.BindConfig) *App {
	handler := apps.NewHandler(appsServer)
	server := apps.NewHttpServer(fmt.Sprintf("%s:%s", cnf.Addr, cnf.Port), handler.GetMuxRouter())
	return &App{
		server: server,
	}
}

func (a *App) Run() error {
	return a.server.Run()
}
