package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
func sessionUser(r *http.Request) *User {
	if u, ok := r.Context().Value("user").(*User); ok {
		return u
	}
	return nil
}

func cookieUser(r *http.Request) *User {
	// previous verions of ident server didn't make use of a domain
	// check for this form of cookie, removing it if found
	// if session, err := sessions.NewCookieStore([]byte(cfg.SessionSecret)).Get(r, cfg.UserCookieKey); err == nil {
	// 	if id, ok := session.Values["id"].(string); ok {
	// 		return NewUser(id)
	// 	}
	// }

	if session, err := sessionStore.Get(r, cfg.UserCookieKey); err == nil {
		if session.Values["id"] != nil {
			if id, ok := session.Values["id"].(string); ok {
				return NewUser(id)
			}
		}
	}
	return nil
}

// check if a user has been provided via access_token param either as a header
// or as request params
func tokenUser(r *http.Request) *User {
	u := NewUser("")
	if r.Header.Get("access_token") != "" {
		u.accessToken = r.Header.Get("access_token")
		return u
	} else if r.FormValue("access_token") != "" {
		u.accessToken = r.FormValue("access_token")
		return u
	} else {
		return nil
	}
}

// attempt to extract & read a session user from a given request.
// if no user is provided, an anonymous user is created
func userFromRequest(db *sql.DB, r *http.Request) (*User, error) {
	var (
		u   *User
		err error
	)

	u = tokenUser(r)
	if u == nil {
		u = cookieUser(r)
	}
	if u == nil {
		u, err = jwtUser(db, r)
		if err != nil {
			// logger.Println(err.Error())
		}
	}

	if u == nil {
		// create anononymous user with ip address
		return &User{
			Username:  getIP(r),
			anonymous: true,
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
	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		return err
	}
	session.Values["id"] = id
	return session.Save(r, w)
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	envelope := r.FormValue("envelope") != "false"
	if u == nil || u.Id == "" {
		w.WriteHeader(http.StatusUnauthorized)
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
			ErrRes(w, NewFmtError(http.StatusBadRequest, "error parsing json: '%s'", err.Error()))
			return
		}
	} else {
		// default to using form values
		login.Username = r.FormValue("username")
		login.Password = r.FormValue("password")
	}

	u, err := AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		ErrRes(w, ErrInvalidUserNamePasswordCombo)
		return
	}

	log.Info("user api login: %s", login.Username)
	if err := setUserSessionCookie(w, r, u.Id); err != nil {
		ErrRes(w, New500Error(err.Error()))
		return
	}

	Res(w, true, u)
}

// logout a user, overwriting their session cookie with ""
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	// previous verions of ident server didn't make use of a domain
	// check for this form of cookie, removing it if found
	// session, err := sessions.NewCookieStore([]byte(cfg.SessionSecret)).Get(r, cfg.UserCookieKey)
	// if err == nil {
	// 	if id, ok := session.Values["id"].(string); ok {
	// 		u := NewUser(id)
	// 		session.Values["id"] = nil
	// 		session.Options.MaxAge = -1
	// 		if err := session.Save(r, w); err != nil {
	// 			log.Infoln(err.Error())
	// 			ErrRes(w, err)
	// 			return
	// 		}
	// 		if err := u.Read(appDB); err == nil {
	// 			log.Info("logout user: %s", u.Username)
	// 		}
	// 		MessageResponse(w, "successfully logged out", nil)
	// 		return
	// 	}
	// }

	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		if session != nil {
			// if we have a session, but have errored for some reason,
			// remove the cookie
			session.Values["id"] = nil
			session.Options.MaxAge = -1
			if err := session.Save(r, w); err != nil {
				log.Infoln(err.Error())
				ErrRes(w, err)
				return
			}
		}
		ErrRes(w, err)
		return
	}

	if id, ok := session.Values["id"].(string); ok {
		u := NewUser(id)
		session.Values["id"] = nil
		session.Options.MaxAge = -1
		if err := session.Save(r, w); err != nil {
			log.Infoln(err.Error())
			ErrRes(w, err)
			return
		}
		if err := u.Read(appDB); err == nil {
			log.Info("logout user: %s", u.Username)
		}
	}

	MessageResponse(w, "successfully logged out", nil)
}
