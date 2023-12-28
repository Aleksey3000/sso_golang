package tests

import (
	ssoV1 "SSO/pkg/proto/sso"
	"SSO/tests/sute"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAppDelete(t *testing.T) {
	ctx, st := sute.New(t)

	const countCases = 10

	var keys []string

	for i := 0; i < countCases; i++ {
		// Add app
		req, err := st.AppsClient.NewApp(ctx, &ssoV1.NewAppRequest{})
		require.NoError(t, err)
		keys = append(keys, string(req.Key))

		// Delete app

		_, err = st.AppsClient.DeleteApp(ctx, &ssoV1.DeleteAppRequest{Key: req.Key})
		require.NoError(t, err)

	}
	if !testUniqueStrings(keys) {
		t.Log(keys)
		t.Error("keys is not Unique")
	}
}

func testUniqueStrings(strings []string) bool {
	for i, s := range strings {
		for j, sub := range strings {
			if j != i && sub == s {
				return false
			}
		}
	}
	return true
}
