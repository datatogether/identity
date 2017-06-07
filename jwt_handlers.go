package main

import (
	"encoding/json"
	"fmt"
	"github.com/archivers-space/identity/jwt"
	"github.com/archivers-space/identity/user"
	"io"
	"net/http"
)

func JwtHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "POST":
		JwtTokenHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

// reads the form values, checks them and creates the token
func JwtTokenHandler(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Username string
		Password string
	}

	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, "error parsing json: '%s'", err.Error()))
			return
		}
	} else {
		// default to using form values
		login.Username = r.Header.Get("username")
		login.Password = r.Header.Get("password")
	}

	u, err := user.AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		ErrRes(w, ErrInvalidUserNamePasswordCombo)
		return
	}

	tokenString, err := jwt.CreateToken(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Info("Token Signing error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, tokenString)
}
