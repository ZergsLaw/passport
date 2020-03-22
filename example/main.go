package main

import (
	"log"

	"github.com/ZergsLaw/login"
	"github.com/ZergsLaw/login/modules/github"
)

// nolint:gochecknoglobals
var (
	githubCfg = login.Config{
		ClientID:     "YOU_APP_ID",
		ClientSecret: "YOU_SECRET_KEY",
		RedirectURI:  "https://oauth2.example.com/",
		Scope:        []string{"YOU SCOPE"},
	}
)

func main() {
	githubClient := github.New(githubCfg)

	login := login.New(
		login.OAuthClient(github.ID, githubClient),
	)

	code := "AUTH_CODE"

	val := github.Account{}
	err := login.Auth(github.ID).Do(code, &val)
	if err != nil {
		log.Fatal(err)
	}
}
