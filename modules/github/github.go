package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ZergsLaw/passport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type (
	client struct {
		cfg oauth2.Config
	}

	account struct {
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Login     string `json:"login"`
	}
)

const ID passport.SocialID = "Github"

// New creates and returns OAuth client.
func New(cfg passport.OAuthConfig) passport.OauthClient {
	return &client{
		cfg: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     github.Endpoint,
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
	const uriUserInfo = "https://api.github.com/user"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return nil, fmt.Errorf("get user from github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, passport.OAuthError(resp.StatusCode, resp.Body)
	}

	githubAccount := account{}
	err = json.NewDecoder(resp.Body).Decode(&githubAccount)
	if err != nil {
		return nil, fmt.Errorf("decode json from github: %w", err)
	}

	return passportAccount(githubAccount), nil
}

func passportAccount(githubAccount account) *passport.Account {
	return &passport.Account{
		ID:     strconv.Itoa(githubAccount.ID),
		Name:   githubAccount.Name,
		Email:  githubAccount.Email,
		Avatar: githubAccount.AvatarURL,
		Login:  githubAccount.Login,
	}
}
