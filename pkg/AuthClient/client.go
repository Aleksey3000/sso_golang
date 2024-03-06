package AuthClient

import (
	ssoV1 "SSO/pkg/proto/sso"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type Client struct {
	appKey           []byte
	authClient       ssoV1.AuthClient
	permissionClient ssoV1.PermissionsClient
}

func New(host string, port string, appKey string) (*Client, error) {
	addr := net.JoinHostPort(host, port)
	cc, err := grpc.DialContext(context.Background(),
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	client := &Client{
		appKey:           []byte(appKey),
		authClient:       ssoV1.NewAuthClient(cc),
		permissionClient: ssoV1.NewPermissionsClient(cc),
	}
	return client, nil

}

func (c *Client) Register(ctx context.Context, login string, password string) error {
	_, err := c.authClient.Register(ctx, &ssoV1.RegisterRequest{
		AppKey:   c.appKey,
		Login:    login,
		Password: password,
	})
	return err
}

func (c *Client) Login(ctx context.Context, login string, password string) (string, error) {
	req, err := c.authClient.Login(ctx, &ssoV1.LoginRequest{
		AppKey:   c.appKey,
		Login:    login,
		Password: password,
	})
	return req.Token, err
}

func (c *Client) DeleteUser(ctx context.Context, login string) error {
	_, err := c.authClient.DeleteUser(ctx, &ssoV1.DeleteUserRequest{
		AppKey: c.appKey,
		Login:  login,
	})
	return err
}

func (c *Client) UpdateLogin(ctx context.Context, login string, newLogin string) error {
	_, err := c.authClient.UpdateLogin(ctx, &ssoV1.UpdateLoginRequest{
		AppKey:   c.appKey,
		Login:    login,
		NewLogin: newLogin,
	})
	return err
}

func (c *Client) ParseToken(ctx context.Context, token string) (login string, err error) {
	req, err := c.authClient.ParseToken(ctx, &ssoV1.ParseTokenRequest{
		AppKey: c.appKey,
		Token:  token,
	})
	return req.Login, err
}

func (c *Client) TestUserOnExist(ctx context.Context, login string) (bool, error) {
	req, err := c.authClient.TestUserOnExist(ctx, &ssoV1.TestUserOnExistRequest{
		AppKey: c.appKey,
		Login:  login,
	})
	return req.Exist, err
}

func (c *Client) GetPermission(ctx context.Context, login string) (int32, error) {
	req, err := c.permissionClient.GetUserPermission(ctx, &ssoV1.GetUserPermissionRequest{
		AppKey: c.appKey,
		Login:  login,
	})
	return req.Permission, err
}

func (c *Client) SetPermission(ctx context.Context, login string, permission int32) error {
	_, err := c.permissionClient.SetUserPermission(ctx, &ssoV1.SetUserPermissionRequest{
		AppKey:     c.appKey,
		Login:      login,
		Permission: permission,
	})
	return err
}
