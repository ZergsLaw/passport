package passport

// Option options for configuring the library core.
type Option func(*Core)

// OAuthClient adds a new oauth client.
func OAuthClient(id SocialID, client OauthClient) Option {
	return func(core *Core) {
		core.oauthClients[id] = client
	}
}
