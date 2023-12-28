package auth

import (
	"SSO/internal/service/auth"
	"SSO/internal/storage/storageErrors"
	ssoV1 "SSO/pkg/proto/sso"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SSOServer struct {
	ssoV1.UnimplementedAuthServer
	ssoV1.UnimplementedAppsServer
	ssoV1.UnimplementedPermissionsServer
	auth Auth
	apps Apps
}

type Apps interface {
	NewApp(ctx context.Context) (key []byte, err error)
	DeleteApp(ctx context.Context, key []byte) (err error)
}

type Auth interface {
	Register(ctx context.Context, appKey []byte, login string, password string) (err error)
	Login(ctx context.Context, appKey []byte, login string, password string) (token string, err error)
}

func RegisterServer(server *grpc.Server, auth Auth, apps Apps) {
	ssoServer := &SSOServer{
		auth: auth,
		apps: apps,
	}
	ssoV1.RegisterAuthServer(server, ssoServer)
	ssoV1.RegisterAppsServer(server, ssoServer)
}

func (s *SSOServer) Register(ctx context.Context, in *ssoV1.RegisterRequest) (*ssoV1.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}

	err := s.auth.Register(ctx, in.AppKey, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, storageErrors.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssoV1.RegisterResponse{}, nil
}

func (s *SSOServer) Login(ctx context.Context, in *ssoV1.LoginRequest) (*ssoV1.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}

	token, err := s.auth.Login(ctx, in.AppKey, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssoV1.LoginResponse{Token: token}, nil
}

func (s *SSOServer) NewApp(ctx context.Context, _ *ssoV1.NewAppRequest) (*ssoV1.NewAppResponse, error) {
	key, err := s.apps.NewApp(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed create app")
	}

	return &ssoV1.NewAppResponse{Key: key}, nil
}

func (s *SSOServer) DeleteApp(ctx context.Context, in *ssoV1.DeleteAppRequest) (*ssoV1.DeleteAppResponse, error) {
	err := s.apps.DeleteApp(ctx, in.Key)
	if err != nil {
		return &ssoV1.DeleteAppResponse{}, status.Error(codes.Internal, "failed delete app")
	}
	return &ssoV1.DeleteAppResponse{}, nil
}
