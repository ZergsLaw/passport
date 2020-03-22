package passport

import (
	"context"
	"io"
)

type (
	// OauthClient the interface to the oauth server.
	// Returns the status and body of the request.
	OauthClient interface {
		Account(ctx context.Context, code string) (status int, body io.ReadCloser, err error)
	}
	// Core the library core, is responsible for managing all modules, and is also a provider to oauth clients.
	Core struct {
		oauthClients map[SocialID]OauthClient
		maxBodySize  int64
	}
	// Config for oauth client.
	Config struct {
		ClientID     string
		ClientSecret string
		RedirectURI  string
		Scope        []string
	}
	// SocialID contains the id of the social network.
	SocialID string
)

const defaultMaxLimitBody = 64 * 1024

// New return instance core.
func New(options ...Option) *Core {
	c := &Core{
		oauthClients: make(map[SocialID]OauthClient),
		maxBodySize:  defaultMaxLimitBody,
	}

	for i := range options {
		options[i](c)
	}

	return c
}

// Auth returns the object for authorization.
func (c *Core) Auth(id SocialID) *Auth {
	client := c.oauthClients[id]

	return &Auth{ctx: context.Background(), client: client, maxBodySize: c.maxBodySize}
}

// AuthWithCtx returns the object to perform authorization with the specified context.
func (c *Core) AuthWithCtx(ctx context.Context, id SocialID) *Auth {
	client := c.oauthClients[id]

	return &Auth{ctx: ctx, client: client, maxBodySize: c.maxBodySize}
}
