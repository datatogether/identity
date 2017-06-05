package main

import (
	"encoding/json"
	"fmt"
	"github.com/archivers-space/identity/user"
	"net/http"
)

func JwtHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	// case "GET":
	// 	GetJwtHandler(w, r)
	// case "PUT":
	// SaveUserHandler(w, r)
	case "POST":
		JwtTokenHandler(w, r)
	// case "DELETE":
	// LogoutHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func JwtPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, cfg.PublicKey)
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
		login.Username = r.FormValue("username")
		login.Password = r.FormValue("password")
	}

	u, err := user.AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		ErrRes(w, ErrInvalidUserNamePasswordCombo)
		return
	}

	tokenString, err := createToken(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Info("Token Signing error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, tokenString)
}
