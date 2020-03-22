package github

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ZergsLaw/login"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type (
	client struct {
		cfg oauth2.Config
	}

	Account struct {
		ID                      int       `json:"id"`
		PrivateGists            int       `json:"private_gists"`
		TotalPrivateRepos       int       `json:"total_private_repos"`
		OwnedPrivateRepos       int       `json:"owned_private_repos"`
		DiskUsage               int       `json:"disk_usage"`
		Collaborators           int       `json:"collaborators"`
		PublicRepos             int       `json:"public_repos"`
		PublicGists             int       `json:"public_gists"`
		Followers               int       `json:"followers"`
		Following               int       `json:"following"`
		Login                   string    `json:"login"`
		NodeID                  string    `json:"node_id"`
		AvatarURL               string    `json:"avatar_url"`
		GravatarID              string    `json:"gravatar_id"`
		URL                     string    `json:"url"`
		HTMLURL                 string    `json:"html_url"`
		FollowersURL            string    `json:"followers_url"`
		FollowingURL            string    `json:"following_url"`
		GistsURL                string    `json:"gists_url"`
		StarredURL              string    `json:"starred_url"`
		SubscriptionsURL        string    `json:"subscriptions_url"`
		OrganizationsURL        string    `json:"organizations_url"`
		ReposURL                string    `json:"repos_url"`
		EventsURL               string    `json:"events_url"`
		ReceivedEventsURL       string    `json:"received_events_url"`
		Type                    string    `json:"type"`
		Name                    string    `json:"name"`
		Company                 string    `json:"company"`
		Blog                    string    `json:"blog"`
		Location                string    `json:"location"`
		Email                   string    `json:"email"`
		Bio                     string    `json:"bio"`
		Hireable                bool      `json:"hireable"`
		SiteAdmin               bool      `json:"site_admin"`
		TwoFactorAuthentication bool      `json:"two_factor_authentication"`
		CreatedAt               time.Time `json:"created_at"`
		UpdatedAt               time.Time `json:"updated_at"`
		Plan                    struct {
			Name          string `json:"name"`
			Space         int    `json:"space"`
			PrivateRepos  int    `json:"private_repos"`
			Collaborators int    `json:"collaborators"`
		} `json:"plan"`
	}
)

const ID login.SocialID = "Github"

// New creates and returns OAuth client.
func New(cfg login.Config) login.OauthClient {
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

func (c *client) Account(ctx context.Context, code string) (status int, body io.ReadCloser, err error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return 0, nil, fmt.Errorf("exchange token: %w", err)
	}

	const uriUserInfo = "https://api.github.com/user"
	resp, err := c.cfg.Client(ctx, token).Get(uriUserInfo)
	if err != nil {
		return 0, nil, fmt.Errorf("get user info from github: %w", err)
	}

	return resp.StatusCode, resp.Body, nil
}
