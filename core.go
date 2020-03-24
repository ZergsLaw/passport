package passport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/oauth2"
)

type (
	// OauthClient the interface to the oauth server.
	// Returns the status and body of the request.
	OauthClient interface {
		Token(ctx context.Context, code string) (*oauth2.Token, error)
		Account(ctx context.Context, token *oauth2.Token) (*Account, error)
	}
	// Core the library core, is responsible for managing all modules, and is also a provider to oauth clients.
	Core struct {
		oauthClients map[SocialID]OauthClient
	}
	// Account contains default user information.
	Account struct {
		ID     string
		Name   string
		Email  string
		Avatar string
	}
	// OAuthConfig for oauth client.
	OAuthConfig struct {
		ClientID     string
		ClientSecret string
		RedirectURI  string
		Scope        []string
	}
	// SocialID contains the id of the social network.
	SocialID string
)

const MaxLimitBody = 64 * 1024

// Errors.
var (
	ErrUnknownOAuthClient = errors.New("unknown oauth client")
)

// New return instance core.
func New(option ...Option) *Core {
	c := &Core{
		oauthClients: make(map[SocialID]OauthClient),
	}

	for i := range option {
		option[i](c)
	}

	return c
}

// Account returns the user account from a specific social network.
// Also returns the network social network token, in case it is necessary
// to call additional methods on behalf of the user.
func (c *Core) Account(ctx context.Context, id SocialID, code string) (*Account, *oauth2.Token, error) {
	client, found := c.oauthClients[id]
	if !found {
		return nil, nil, ErrUnknownOAuthClient
	}

	token, err := client.Token(ctx, code)
	if err != nil {
		return nil, nil, err
	}

	account, err := client.Account(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	return account, token, nil
}

// OAuthErr error from OAuth server.
type OAuthErr struct {
	Code int
	Body json.RawMessage
}

func (err OAuthErr) Error() string {
	return fmt.Sprintf("code: %d, text: %s", err.Code, err.Body)
}

func OAuthError(status int, r io.Reader) error {
	body, err := ioutil.ReadAll(io.LimitReader(r, MaxLimitBody))
	if err != nil {
		return fmt.Errorf("passport: read body error: %w", err)
	}

	return OAuthErr{Code: status, Body: body}
}
