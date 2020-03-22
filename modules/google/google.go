package google

import (
	"context"
	"fmt"
	"io"

	"github.com/ZergsLaw/passport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	client struct {
		cfg oauth2.Config
	}

	Account struct {
		Email         string `json:"email"`
		FamilyName    string `json:"family_name"`
		GivenName     string `json:"given_name"`
		ID            string `json:"id"`
		Locale        string `json:"locale"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}
)

const ID passport.SocialID = "Google"

// New creates and returns OAuth client.
func New(cfg passport.Config) passport.OauthClient {
	return &client{
		cfg: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  cfg.RedirectURI,
			Scopes:       cfg.Scope,
		},
	}
}

func (c *client) Account(ctx context.Context, code string) (status int, body io.ReadCloser, err error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return 0, nil, fmt.Errorf("exchange token: %w", err)
	}

	const uriUserInfo = "https://www.googleapis.com/oauth2/v2/userinfo?alt=json"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return 0, nil, fmt.Errorf("get user info from google: %w", err)
	}

	return resp.StatusCode, resp.Body, nil
}
