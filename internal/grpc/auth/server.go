package auth

import (
	"SSO/internal/service/auth"
	"SSO/internal/storage/storageErrors"
	ssov1 "SSO/pkg/proto/sso"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Register(ctx context.Context, appKey []byte, login string, password string) (err error)
	Login(ctx context.Context, appKey []byte, login string, password string) (token string, err error)
}

func RegisterServer(server *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(server, &Server{
		auth: auth,
	})
}

func (s *Server) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
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

		return nil, status.Error(codes.Internal, err.Error()) //"failed to register user"
	}

	return &ssov1.RegisterResponse{}, nil
}

func (s *Server) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
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

	return &ssov1.LoginResponse{Token: token}, nil
}
