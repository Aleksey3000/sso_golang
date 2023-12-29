package tests

import (
	ssoV1 "SSO/pkg/proto/sso"
	"SSO/tests/sute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetGetPermission(t *testing.T) {
	ctx, st := sute.New(t)

	var (
		userId int64 = 1
		perm   int32 = 1
	)

	_, err := st.PermClient.SetUserPermission(ctx, &ssoV1.SetUserPermissionRequest{
		UserId:     userId,
		Permission: perm,
		AppKey:     appKey,
	})
	require.NoError(t, err)

	resp, err := st.PermClient.GetUserPermission(ctx, &ssoV1.GetUserPermissionRequest{
		UserId: userId,
		AppKey: appKey,
	})
	require.NoError(t, err)
	assert.Equal(t, perm, resp.Permission)
}
