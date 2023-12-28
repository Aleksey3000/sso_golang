package grpc

import (
	"SSO/internal/config"
	"SSO/internal/grpc/auth"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	l          *slog.Logger
	grpcServer *grpc.Server
	bindCnf    *config.BindConfig
}

// protoc -I proto proto/sso/sso.proto --go_out=./pkg/proto/ --go_opt=paths=source_relative --go-grpc_out=./pkg/proto/ --go-grpc_opt=paths=source_relative

func New(l *slog.Logger, authService auth.Auth, cnf *config.BindConfig) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			l.Error("recovered panic ", slog.Any("panic: ", p))
			return status.Error(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLog(l), loggingOpts...),
	))

	auth.RegisterServer(grpcServer, authService)

	return &App{
		l:          l,
		grpcServer: grpcServer,
		bindCnf:    cnf,
	}
}

func (a *App) Run() error {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.bindCnf.Addr, a.bindCnf.Port))
	if err != nil {
		return err
	}
	if err := a.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}

func interceptorLog(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}
