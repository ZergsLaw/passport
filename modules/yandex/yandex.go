package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ZergsLaw/passport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

type (
	client struct {
		cfg oauth2.Config
	}

	account struct {
		ID           string `json:"id"`
		DisplayName  string `json:"display_name"`
		DefaultEmail string `json:"default_email"`
		Login        string `json:"login"`
	}
)

const ID passport.SocialID = "Yandex"

// New creates and returns OAuth client.
func New(cfg passport.OAuthConfig) passport.OauthClient {
	return &client{
		cfg: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     yandex.Endpoint,
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
	const uriUserInfo = "https://login.yandex.ru/info?format=json"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return nil, fmt.Errorf("get user from github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, passport.OAuthError(resp.StatusCode, resp.Body)
	}

	yandexAccount := account{}
	err = json.NewDecoder(resp.Body).Decode(&yandexAccount)
	if err != nil {
		return nil, fmt.Errorf("decode json from github: %w", err)
	}

	return passportAccount(yandexAccount), nil
}

func passportAccount(yandexAccount account) *passport.Account {
	return &passport.Account{
		ID:    yandexAccount.ID,
		Name:  yandexAccount.DisplayName,
		Email: yandexAccount.DefaultEmail,
		Login: yandexAccount.Login,
	}
}
