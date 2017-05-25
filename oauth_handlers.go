package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
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
				logger.Println(err.Error())
				ErrRes(w, err)
				return
			}

			if info["login"] == nil {
				logger.Println("no user found")
				ErrRes(w, fmt.Errorf("no user found"))
				return
			}

			perm, err := g.RepoPermission(r.FormValue("owner"), r.FormValue("repo"), info["login"].(string))
			if err != nil {
				logger.Println(err.Error())
				ErrRes(w, err)
				return
			}

			Res(w, true, map[string]string{"permission": perm})
			return
		}
	}

	err = NewFmtError(http.StatusUnauthorized, "this user hasn't enabled github for their account")
	logger.Println(err.Error())
	ErrRes(w, err)
}

// redirect user to github auth url
func GithubOauthHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = cfg.UrlRoot
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(redirect))
	url := githubOAuth.AuthCodeURL(b64, oauth2.AccessTypeOffline)
	// logger.Println("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handle Oauth response from github
// TODO - refactor ASAP.
func GithubOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user := sessionUser(r)
	ctx := r.Context()

	redirectBytes, err := base64.StdEncoding.DecodeString(r.FormValue("state"))
	if err != nil {
		logger.Println(err.Error())
		ErrRes(w, fmt.Errorf("bad response value: %s", err.Error()))
		return
	}
	redirect := string(redirectBytes)

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

	if err := t.Read(appDB); err == nil {
		if err := setUserSessionCookie(w, r, t.User.Id); err != nil {
			ErrRes(w, fmt.Errorf("error setting session cookie: %s", err.Error()))
			return
		}

		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}
	} else if user.anonymous {
		svc, err := t.UserService()
		if err != nil {
			logger.Println(err.Error())
			ErrRes(w, err)
			return
		}
		u, err := svc.ExtractUser()
		if err != nil {
			logger.Println(err.Error())
			ErrRes(w, err)
			return
		}
		t.User = u

		emailUser := &User{Email: u.Email}
		if err := emailUser.Read(appDB); err == nil {
			// if we have a matching email, connect the two accounts
			t.User = emailUser
		} else if err := t.User.Save(appDB); err != nil {
			// create a new user that matches
			// TODO - better username collision handling
			logger.Println(err)
			if err == ErrUsernameTaken {
				for i := 1; i < 1000; i++ {
					t.User.Username = fmt.Sprintf("%s_%d", t.User.Username, i)
					if err := t.User.Save(appDB); err == nil {
						break
					} else if err != ErrUsernameTaken {
						logger.Println(err.Error())
						ErrRes(w, err)
						return
					}
				}
			} else {
				logger.Println(err.Error())
				ErrRes(w, err)
				return
			}
		}
	}

	if err := t.Save(appDB); err != nil {
		logger.Println(err.Error())
		ErrRes(w, err)
		return
	}

	if err := setUserSessionCookie(w, r, t.User.Id); err != nil {
		ErrRes(w, fmt.Errorf("error setting session cookie: %s", err.Error()))
		return
	}

	if redirect != "" {
		redirect, err := ValidUrlString(redirect)
		if err == nil {
			logger.Println("redirecting to", redirect)
			http.Redirect(w, r, redirect, http.StatusFound)
		}
		return
	}

	http.Redirect(w, r, cfg.FrontendUrl, http.StatusFound)

	// client := githubOAuth.Client(ctx, tok)
	// res, err := client.Get("https://api.github.com/repos/edgi-govdata-archiving/archivers.space/collaborators")
	// if err != nil {
	// 	logger.Println(err.Error())
	// 	ErrRes(w, err)
	// 	return
	// }
	// defer res.Body.Close()
	// data, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	logger.Println(err.Error())
	// 	ErrRes(w, err)
	// 	return
	// }
	// w.Write(data)
}
