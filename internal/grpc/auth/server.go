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
	ssoV1.UnimplementedPermissionsServer

	auth        Auth
	permissions Permissions

	apps Apps
}

type Apps interface {
	TestOnExist(ctx context.Context, key []byte) bool
}

type Auth interface {
	Register(ctx context.Context, appKey []byte, login string, password string) (err error)
	Login(ctx context.Context, appKey []byte, login string, password string) (token string, err error)
	DeleteUser(ctx context.Context, appKey []byte, login string) (err error)
	UpdateLogin(ctx context.Context, appKey []byte, login string, newLogin string) error
	ChangePassword(ctx context.Context, appKey []byte, login string, newPass string) error
	TestOnExist(ctx context.Context, appKey []byte, login string) bool
	GetUserId(ctx context.Context, appKey []byte, login string) (int64, error)
	ParseToken(ctx context.Context, appKey []byte, token string) (string, error)
}

type Permissions interface {
	SetUserPermission(ctx context.Context, userId int64, permission int32) (err error)
	GetUserPermission(ctx context.Context, userId int64) (permission int32, err error)
}

func RegisterServer(server *grpc.Server, auth Auth, apps Apps, permissions Permissions) {
	ssoServer := &SSOServer{
		auth:        auth,
		permissions: permissions,
		apps:        apps,
	}
	ssoV1.RegisterAuthServer(server, ssoServer)
	ssoV1.RegisterPermissionsServer(server, ssoServer)
}

var ErrNilRequest = errors.New("nil request")

func (s *SSOServer) Register(ctx context.Context, in *ssoV1.RegisterRequest) (*ssoV1.RegisterResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
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
	if in == nil {
		return nil, ErrNilRequest
	}
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

func (s *SSOServer) DeleteUser(ctx context.Context, in *ssoV1.DeleteUserRequest) (*ssoV1.DeleteUserResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	err := s.auth.DeleteUser(ctx, in.AppKey, in.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed delete user")
	}
	return &ssoV1.DeleteUserResponse{}, err
}

func (s *SSOServer) UpdateLogin(ctx context.Context, in *ssoV1.UpdateLoginRequest) (*ssoV1.UpdateLoginResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if in.NewLogin == "" {
		return nil, status.Error(codes.InvalidArgument, "new login is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	err := s.auth.UpdateLogin(ctx, in.AppKey, in.Login, in.NewLogin)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed update login")
	}
	return &ssoV1.UpdateLoginResponse{}, err
}

func (s *SSOServer) ChangePassword(ctx context.Context, in *ssoV1.ChangePasswordRequest) (*ssoV1.ChangePasswordResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	if in.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}
	if err := s.auth.ChangePassword(ctx, in.AppKey, in.Login, in.NewPassword); err != nil {
		return nil, status.Error(codes.Internal, "failed change password")
	}

	return &ssoV1.ChangePasswordResponse{}, nil
}

func (s *SSOServer) TestUserOnExist(ctx context.Context, in *ssoV1.TestUserOnExistRequest) (*ssoV1.TestUserOnExistResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}

	exist := s.auth.TestOnExist(ctx, in.AppKey, in.Login)
	return &ssoV1.TestUserOnExistResponse{Exist: exist}, nil
}

func (s *SSOServer) ParseToken(ctx context.Context, in *ssoV1.ParseTokenRequest) (*ssoV1.ParseTokenResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	login, err := s.auth.ParseToken(ctx, in.AppKey, in.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "err in parse token: %s", err.Error())
	}
	return &ssoV1.ParseTokenResponse{Login: login}, err
}

func (s *SSOServer) GetUserPermission(ctx context.Context, in *ssoV1.GetUserPermissionRequest) (*ssoV1.GetUserPermissionResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if !s.apps.TestOnExist(ctx, in.AppKey) {
		return nil, status.Error(codes.FailedPrecondition, "app not found")
	}

	id, err := s.auth.GetUserId(ctx, in.AppKey, in.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "user not found")
	}

	perm, err := s.permissions.GetUserPermission(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed get user permission")
	}
	return &ssoV1.GetUserPermissionResponse{Permission: perm}, nil
}

func (s *SSOServer) SetUserPermission(ctx context.Context, in *ssoV1.SetUserPermissionRequest) (*ssoV1.SetUserPermissionResponse, error) {
	if in == nil {
		return nil, ErrNilRequest
	}
	if len(in.AppKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "app key is required")
	}
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if !s.apps.TestOnExist(ctx, in.AppKey) {
		return nil, status.Error(codes.FailedPrecondition, "app not found")
	}

	id, err := s.auth.GetUserId(ctx, in.AppKey, in.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "user not found")
	}

	if err := s.permissions.SetUserPermission(ctx, id, in.Permission); err != nil {
		return nil, status.Error(codes.Internal, "failed set permission")
	}
	return &ssoV1.SetUserPermissionResponse{}, nil
}
