package main

import (
	"encoding/json"
	"net/http"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		ListUsersHandler(w, r)
	case "POST":
		CreateUserHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		SingleUserHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func SingleUserHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"
	p := &UsersGetParams{
		Subject: &User{
			Id:          r.FormValue("id"),
			Username:    r.FormValue("username"),
			accessToken: r.FormValue("access_token"),
		},
	}
	res := &User{}
	if err := new(Users).Get(p, res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, envelope, res)
}

// list users or get a single user if supplied with a "username" formValue
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"
	username := r.FormValue("username")
	id := r.FormValue("id")

	if username != "" || id != "" {
		SingleUserHandler(w, r)
		return
	} else {
		p := &UsersListParams{
			User: sessionUser(r),
			Page: PageFromRequest(r),
		}
		res := []*User{}
		if err := new(Users).List(p, &res); err != nil {
			ErrRes(w, err)
			return
		}
		Res(w, envelope, res)
	}

	// ExecRequest(w, envelope, req)
}

// Create a user from the api, feed password in as a query param
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// sess := sessionUser(r)
	u := NewUser("")

	if isJsonRequest(r) {
		params := struct {
			Username string
			Email    string
			Password string
		}{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, "error decoding json: %s", err.Error()))
			return
		}
		u = &User{
			Username: params.Username,
			Email:    params.Email,
			password: params.Password,
		}
	} else {
		// default to form data values
		u.Username = r.FormValue("username")
		u.Email = r.FormValue("email")
		u.password = r.FormValue("password")
	}

	p := &UsersCreateParams{
		User: u,
	}

	res := &User{}
	if err := new(Users).Create(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	// log new user in
	if err := setUserSessionCookie(w, r, u.Id); err != nil {
		ErrRes(w, New500Error(err.Error()))
		return
	}

	Res(w, true, res)
}

// confirm a user's email address
// func UserConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
// 	u := NewUser(stringParam(ctx, "id"))
// 	u.emailConfirmed = true
// 	log.Info(u)
// 	if err := u.Save(appDB); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	if err := AddFlashMessage(w, r, "Thanks! Your email address has been confirmed."); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	http.Redirect(w, r, u.Path(), http.StatusSeeOther)
// }

// SaveUserHandler updates a user
func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	u := &User{}
	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, err.Error()))
			return
		}
	} else {
		// TODO - fill user out from form values
	}

	p := &UsersSaveParams{
		User:    sessionUser(r),
		Subject: u,
	}

	res := &User{}
	if err := new(Users).Save(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

func UsersSearchHandler(w http.ResponseWriter, r *http.Request) {
	p := &UsersSearchParams{
		User:  sessionUser(r),
		Query: r.FormValue("q"),
		Page:  PageFromRequest(r),
	}

	res := []*User{}
	if err := new(Users).Search(p, &res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, true, res)
}

// delete a user
// func DeleteCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
// 	u := sessionUser(ctx)
// 	if err := u.Delete(appDB); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	session, err := sessionStore.Get(r, "qri.io.user")
// 	if err != nil {
// 		ErrRes(w, err)
// 		return
// 	}
// 	session.Values["id"] = nil
// 	if err := session.Save(r, w); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	AddFlashMessage(w, r, "You've successfully deleted your account.")
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
