package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

func SessionUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	tokens, err := sessionUser(r).OauthTokens(appDB)
	if err != nil {
		ErrRes(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logger.Println(err.Error())
	}
}

func GithubRepoAccessHandler(w http.ResponseWriter, r *http.Request) {
	tokens, err := sessionUser(r).OauthTokens(appDB)
	if err != nil {
		ErrRes(w, err)
		return
	}

	for _, t := range tokens {
		if t.Service == OauthServiceGithub {
			g := NewGithub(t.token)
			info, err := g.CurrentUserInfo()
			if err != nil {
				ErrRes(w, err)
				return
			}

			if info["login"] == nil {
				ErrRes(w, fmt.Errorf("no user found"))
				return
			}

			perm, err := g.RepoPermission(r.FormValue("owner"), r.FormValue("repo"), info["login"].(string))
			if err != nil {
				ErrRes(w, err)
				return
			}

			Res(w, false, map[string]string{"permission": perm})
			return
		}
	}

	ErrRes(w, fmt.Errorf("this user hasn't enabled github for their account"))
}

func GithubOauthHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := githubOAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// logger.Println("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// TODO - like woah
func GithubOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user := sessionUser(r)
	ctx := r.Context()

	// if state := r.FormValue("state"); state != oauth2.AccessTypeOffline {
	//  }

	service := OauthServiceGithub
	// if service != "github" {
	// 	logger.Fatal(err)
	// }

	code := r.FormValue("code")
	tok, err := githubOAuth.Exchange(ctx, code)
	if err != nil {
		logger.Fatal(err)
	}

	t := &UserOauthToken{
		User:    user,
		Service: service,
		token:   tok,
	}

	if user.anonymous {
		// if err := user.Save(appDB); err != nil {
		// 	logger.Println(err.Error())
		// }
		ser, err := t.UserService()
		if err != nil {
			logger.Println(err.Error())
			return
		}
		u, err := ser.ExtractUser()
		if err != nil {
			logger.Println(err.Error())
			return
		}
		t.User = u
		if err := t.User.Save(appDB); err != nil {
			logger.Println(err.Error())
			return
		}
	}

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
