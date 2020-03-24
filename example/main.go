package main

import (
	"context"
	"log"
	"time"

	"github.com/ZergsLaw/passport"
	"github.com/ZergsLaw/passport/modules/github"
)

var (
	githubCfg = passport.OAuthConfig{
		ClientID:     "GITHUB_APP_ID",
		ClientSecret: "GITHUB_SECRET_KEY",
		RedirectURI:  "https://oauth2.example.com/",
		Scope:        []string{"YOU SCOPE"},
	}
)

func main() {
	githubCfg = passport.OAuthConfig{
		ClientID:     "GITHUB_APP_ID",
		ClientSecret: "GITHUB_SECRET_KEY",
		RedirectURI:  "https://oauth2.example.com/",
		Scope:        []string{"YOU SCOPE"},
	}

	login := passport.New(
		github.OAuth(githubCfg),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	code := "AUTH_CODE"

	account, token, err := login.Account(ctx, github.ID, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("account:", account)
	log.Println("token type:", token.TokenType)
}
