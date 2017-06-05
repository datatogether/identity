package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/archivers-space/identity/users"
	"net"
	"net/http"
	"strings"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		GetSessionHandler(w, r)
	case "PUT":
		SaveUserHandler(w, r)
	case "POST":
		LoginHandler(w, r)
	case "DELETE":
		LogoutHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

// grab the current serialized sesion user, returns nil if no user present
func sessionUser(r *http.Request) *users.User {
	if u, ok := r.Context().Value("user").(*users.User); ok {
		return u
	}
	return nil
}

func cookieUser(r *http.Request) *users.User {
	if session, err := sessionStore.Get(r, cfg.UserCookieKey); err == nil {
		if session.Values["id"] != nil {
			if id, ok := session.Values["id"].(string); ok {
				return users.NewUser(id)
			}
		}
	} else {
		log.Infoln(err.Error())
	}
	return nil
}

// check if a user has been provided via access_token param either as a header
// or as request params
func tokenUser(r *http.Request) *users.User {
	u := users.NewUser("")
	if r.Header.Get("access_token") != "" {
		u = users.NewAccessTokenUser(r.Header.Get("access_token"))
		return u
	} else if r.FormValue("access_token") != "" {
		u = users.NewAccessTokenUser(r.FormValue("access_token"))
		return u
	} else {
		return nil
	}
}

// attempt to extract & read a session user from a given request.
// if no user is provided, an anonymous user is created
func userFromRequest(db *sql.DB, r *http.Request) (*users.User, error) {
	var (
		u   *users.User
		err error
	)

	u = tokenUser(r)
	if u == nil {
		u, err = jwtUser(db, r)
		if err != nil {
			// logger.Println(err.Error())
		}
	}
	if u == nil {
		u = cookieUser(r)
	}
	if u == nil {
		// create anononymous user with ip address
		return &users.User{
			Username:  getIP(r),
			Anonymous: true,
		}, nil
	}

	if err := u.Read(db); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return u, nil
}

func getIP(r *http.Request) string {
	remoteAddr := r.Header.Get("x-forwarded-for")
	if remoteAddr != "" {
		return strings.TrimSpace(strings.Split(remoteAddr, ",")[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

// set a user's session cookie so we can track them
func setUserSessionCookie(w http.ResponseWriter, r *http.Request, id string) error {
	log.Infoln("set user cookie:", id)
	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		log.Infoln("setUserSessionCookie error", err.Error())
		if session != nil {
			// if we still get a session object back
			// clear the cookie, b/c this one clearly doesn't work
			session.Options.MaxAge = -1
		}
	}
	session.Values["id"] = id
	return session.Save(r, w)
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	envelope := r.FormValue("envelope") != "false"
	if u == nil || u.Id == "" {
		ErrRes(w, NewFmtError(http.StatusUnauthorized, "unauthorized"))
		return
	} else {
		Res(w, envelope, u)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Username string
		Password string
	}

	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			e := NewFmtError(http.StatusBadRequest, "error parsing json: '%s'", err.Error())
			log.Infoln(e)
			ErrRes(w, e)
			return
		}
	} else {
		// default to using form values
		login.Username = r.FormValue("username")
		login.Password = r.FormValue("password")
	}

	u, err := users.AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		log.Infoln(ErrInvalidUserNamePasswordCombo)
		ErrRes(w, ErrInvalidUserNamePasswordCombo)
		return
	}

	if err := setUserSessionCookie(w, r, u.Id); err != nil {
		log.Infoln(err)
		// ErrRes(w, New500Error(err.Error()))
		// return
	}

	log.Infof("user api login: %s", login.Username)
	Res(w, true, u)
}

// logout a user, overwriting their session cookie with ""
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		log.Infof("session get user error", err.Error())
	} else {
		if id, ok := session.Values["id"].(string); ok {
			u := users.NewUser(id)
			if err := u.Read(appDB); err == nil {
				log.Info("logout user: %s", u.Username)
			}
		}
	}

	// regardless of what happens in relation to errors
	// if we have a session object, clear it.
	if session != nil {
		session.Values["id"] = nil
		session.Options.MaxAge = -1
	}
	if err := session.Save(r, w); err != nil {
		log.Infoln(err.Error())
		ErrRes(w, err)
		return
	}

	MessageResponse(w, "successfully logged out", nil)
}
