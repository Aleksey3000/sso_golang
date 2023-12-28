package app

import (
	GrpcApp "SSO/internal/app/grpc"
	"SSO/internal/config"
	"SSO/internal/service/auth"
	"SSO/internal/storage"
	"log/slog"
)

type App struct {
	GRPCApp *GrpcApp.App
}

func New(l *slog.Logger, cnf *config.Config) *App {
	s, err := storage.New(&cnf.DBConfig)
	if err != nil {
		panic(err)
	}

	authService := auth.New(l, s.UserStorage, s.AppStorage, cnf.TokenTTL)
	grpcApp := GrpcApp.New(l, authService, &cnf.BindConfig)
	return &App{
		GRPCApp: grpcApp,
	}
}
