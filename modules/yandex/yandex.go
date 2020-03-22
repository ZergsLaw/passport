package yandex

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ZergsLaw/login"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

type (
	client struct {
		cfg oauth2.Config
	}

	Account struct {
		ID              string    `json:"id"`
		FirstName       string    `json:"first_name"`
		LastName        string    `json:"last_name"`
		DisplayName     string    `json:"display_name"`
		Emails          []string  `json:"emails"`
		DefaultEmail    string    `json:"default_email"`
		RealName        string    `json:"real_name"`
		IsAvatarEmpty   bool      `json:"is_avatar_empty"`
		BirthDay        time.Time `json:"birth_day"`
		DefaultAvatarID string    `json:"default_avatar_id"`
		Login           string    `json:"login"`
		Sex             string    `json:"sex"`
	}
)

const ID login.SocialID = "YA"

// New creates and returns OAuth client.
func New(cfg login.Config) login.OauthClient {
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

func (c *client) Account(ctx context.Context, code string) (status int, body io.ReadCloser, err error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return 0, nil, fmt.Errorf("exchange token: %w", err)
	}

	const uriUserInfo = "https://login.yandex.ru/info?format=json"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return 0, nil, fmt.Errorf("get user info from yandex: %w", err)
	}

	return resp.StatusCode, resp.Body, nil
}
