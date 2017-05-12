package main

import (
	"golang.org/x/oauth2"
)

const (
	OauthServiceGithub = "github"
)

type OauthUserService interface {
	CurrentUserInfo() (map[string]interface{}, error)
	ExtractUser() (*User, error)
}

var (
	githubOAuth *oauth2.Config
)

func initOauth() {
	githubOAuth = &oauth2.Config{
		ClientID:     cfg.GithubAppId,
		ClientSecret: cfg.GithubAppSecret,
		// RedirectURL:  "https://ident.archivers.space/oauth/callback",
		Scopes: []string{"user", "repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}
