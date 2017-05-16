package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

// Return all Token types for the currently logged-in user
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

// Check currently logged-in user's access to a provided github repo
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

			Res(w, true, map[string]string{"permission": perm})
			return
		}
	}

	ErrRes(w, NewFmtError(http.StatusUnauthorized, "this user hasn't enabled github for their account"))
}

// redirect user to github auth url
func GithubOauthHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := githubOAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// logger.Println("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handle Oauth response from github
func GithubOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user := sessionUser(r)
	ctx := r.Context()

	// TODO - MITM check
	// if state := r.FormValue("state"); state != oauth2.AccessTypeOffline {
	//  }

	code := r.FormValue("code")
	tok, err := githubOAuth.Exchange(ctx, code)
	if err != nil {
		logger.Fatal(err)
	}

	t := &UserOauthToken{
		User:    user,
		Service: OauthServiceGithub,
		token:   tok,
	}

	if user.anonymous {
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
		logger.Println(err.Error())
		ErrRes(w, err)
		return
	}

	client := githubOAuth.Client(ctx, tok)
	res, err := client.Get("https://api.github.com/repos/edgi-govdata-archiving/archivers.space/collaborators")
	if err != nil {
		logger.Println(err.Error())
		ErrRes(w, err)
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println(err.Error())
		ErrRes(w, err)
		return
	}

	w.Write(data)
}
