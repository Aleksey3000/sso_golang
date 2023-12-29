package sute

import (
	"SSO/internal/config"
	ssoV1 "SSO/pkg/proto/sso"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

const configPath = "C:\\Users\\79212\\GolandProjects\\SSO\\config\\config.yaml"

type Suite struct {
	*testing.T                  // Потребуется для вызова методов *testing.T
	Cnf        *config.Config   // Конфигурация приложения
	AuthClient ssoV1.AuthClient // Клиент для взаимодействия с gRPC-сервером Auth
	AppsClient ssoV1.AppsClient
	PermClient ssoV1.PermissionsClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	cnf, err := config.GetConfig(configPath)
	t.Logf("%+v", cnf)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Hour*10)

	t.Cleanup(
		func() {
			t.Helper()
			cancelCtx()
		})

	grpcAddress := net.JoinHostPort(cnf.BindConfig.Addr, cnf.BindConfig.Port)

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatal(err)
	}

	authClient := ssoV1.NewAuthClient(cc)
	appsClient := ssoV1.NewAppsClient(cc)
	permClient := ssoV1.NewPermissionsClient(cc)
	return ctx, &Suite{
		T:          t,
		Cnf:        cnf,
		AuthClient: authClient,
		AppsClient: appsClient,
		PermClient: permClient,
	}
}
