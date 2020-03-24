package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZergsLaw/passport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	client struct {
		cfg oauth2.Config
	}

	account struct {
		Email   string `json:"email"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
)

const ID passport.SocialID = "Google"

// New creates and returns OAuth client.
func New(cfg passport.OAuthConfig) passport.OauthClient {
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

func OAuth(cfg passport.OAuthConfig) passport.Option {
	return passport.OAuthClient(ID, New(cfg))
}

func (c *client) Token(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %w", err)
	}
	return token, nil
}

func (c *client) Account(ctx context.Context, token *oauth2.Token) (*passport.Account, error) {
	const uriUserInfo = "https://www.googleapis.com/oauth2/v2/userinfo?alt=json"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return nil, fmt.Errorf("get user from github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, passport.OAuthError(resp.StatusCode, resp.Body)
	}

	googleAccount := account{}
	err = json.NewDecoder(resp.Body).Decode(&googleAccount)
	if err != nil {
		return nil, fmt.Errorf("decode json from github: %w", err)
	}

	return passportAccount(googleAccount), nil
}

func passportAccount(googleAccount account) *passport.Account {
	return &passport.Account{
		ID:     googleAccount.ID,
		Name:   googleAccount.Name,
		Email:  googleAccount.Email,
		Avatar: googleAccount.Picture,
	}
}
