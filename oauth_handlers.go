package main

import (
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

func HandleGithubOauth(w http.ResponseWriter, r *http.Request) {
	githubOAuth = &oauth2.Config{
		ClientID:     cfg.GithubAppId,
		ClientSecret: cfg.GithubAppSecret,
		// RedirectURL:  "https://ident.archivers.space/oauth.callback",
		Scopes: []string{"user", "repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := githubOAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// logger.Println("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// TODO - like woah
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	user := sessionUser(r)
	ctx := r.Context()

	// if state := r.FormValue("state"); state != oauth2.AccessTypeOffline {

	//  }

	service := "github"
	// if service != "github" {
	// 	logger.Fatal(err)
	// }

	code := r.FormValue("code")
	tok, err := githubOAuth.Exchange(ctx, code)
	if err != nil {
		logger.Fatal(err)
	}

	if user.anonymous {
		if err := user.Save(appDB); err != nil {
			logger.Println(err.Error())
		}
	}

	t := &OauthToken{User: user, Service: service}
	if err := t.Save(appDB); err != nil {
		logger.Fatal(err)
	}

	client := githubOAuth.Client(ctx, tok)
	res, err := client.Get("https://api.github.com/repos/edgi-govdata-archiving/archivers.space/collaborators")
	if err != nil {
		logger.Fatal(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatal(err)
	}

	w.Write(data)
}
