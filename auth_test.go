package login_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ZergsLaw/login"
	"github.com/ZergsLaw/login/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type mockAccount struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Blog     string `json:"blog"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Hireable bool   `json:"hireable"`
	Bio      string `json:"bio"`
}

const (
	jsonAccount = `
{
  "name": "monalisa octocat",
  "email": "octocat@github.com",
  "blog": "https://github.com/blog",
  "company": "GitHub",
  "location": "San Francisco",
  "hireable": true,
  "bio": "There once..."
}
`
	jsonErr = `{"err_description":"not found"}`

	oAuthMockID    login.SocialID = "mock"
	unknownOAuthID login.SocialID = "unknown"
	oAuthCode                     = "123456"
)

func TestAuth_Do(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock.NewOauthClient(ctrl)
	core := login.New(login.OAuthClient(oAuthMockID, mockAuth))
	ctx := context.Background()
	assert := require.New(t)
	anyErr := errors.New("any error")

	expectRes := &mockAccount{}
	err := json.Unmarshal([]byte(jsonAccount), expectRes)
	assert.NoError(err)
	oAuthError := &login.OAuthErr{}
	err = json.Unmarshal([]byte(jsonErr), &oAuthError.Value)
	assert.NoError(err)

	successBody := ioutil.NopCloser(bytes.NewReader([]byte(jsonAccount)))
	mockAuth.EXPECT().Account(ctx, oAuthCode).Return(http.StatusOK, successBody, nil)
	mockAuth.EXPECT().Account(ctx, oAuthCode).Return(0, nil, anyErr)

	errBody := ioutil.NopCloser(bytes.NewReader([]byte(jsonErr)))
	mockAuth.EXPECT().Account(ctx, oAuthCode).Return(http.StatusNotFound, errBody, nil)

	testCases := []struct {
		name string
		id   login.SocialID

		expectResult *mockAccount
		expectErr    error
	}{
		{"success", oAuthMockID, expectRes, nil},
		{"unknown_oauth_client", unknownOAuthID, nil, login.ErrUnknownOAuthClient},
		{"any_error_from_oauth_server", oAuthMockID, nil, anyErr},
		{"status_not_ok", oAuthMockID, nil, oAuthError},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			res := &mockAccount{}
			err := core.AuthWithCtx(ctx, tc.id).Do(oAuthCode, res)
			if tc.expectErr == nil {
				assert.NoError(err)
				assert.Equal(tc.expectResult, res)
			} else {
				unwrapErr := errors.Unwrap(err)
				if unwrapErr == nil {
					unwrapErr = err
				}
				assert.Equal(tc.expectErr, unwrapErr)
				assert.Zero(*res)
			}
		})
	}
}
