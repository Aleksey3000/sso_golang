package app

import (
	GrpcApp "SSO/internal/app/grpc"
	HttpApp "SSO/internal/app/http"
	"SSO/internal/config"
	"SSO/internal/service/apps"
	"SSO/internal/service/auth"
	"SSO/internal/service/permissions"
	"SSO/internal/storage"
	"log/slog"
)

type App struct {
	GRPCApp *GrpcApp.App
	HTTPApp *HttpApp.App
}

func New(l *slog.Logger, cnf *config.Config) *App {
	s, err := storage.New(&cnf.DBConfig)
	if err != nil {
		panic(err)
	}

	permService := permissions.New(l, s.PermissionsStorage)
	authService := auth.New(l, s.UserStorage, s.AppStorage, permService, cnf.TokenTTL)
	appsService := apps.New(l, s.AppStorage)

	grpcApp := GrpcApp.New(l, authService, appsService, permService, &cnf.GRPCBindConfig)
	httpApp := HttpApp.NewHttpApp(appsService, &cnf.HttpBindConfig)

	return &App{
		GRPCApp: grpcApp,
		HTTPApp: httpApp,
	}
}
