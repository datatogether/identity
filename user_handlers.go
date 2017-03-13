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

// list users or get a single user if supplied with a "username" formValue
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	username := r.FormValue("username")
	id := r.FormValue("id")

	if username != "" || id != "" {
		req = &UserRequest{
			Interface: httpApiInterface,
			User:      sessionUser(r),
			Subject: &User{
				Id:       id,
				Username: username,
			},
		}
	} else {
		req = &UsersRequest{
			Interface: httpApiInterface,
			User:      sessionUser(r),
			Page:      PageFromRequest(r),
		}
	}

	ExecRequest(w, req)
}

// Create a user from the api, feed password in as a query param
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// sess := sessionUser(r)
	// u := NewUser("")
	params := struct {
		Username string
		Email    string
		Password string
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ErrRes(w, NewFmtError(http.StatusBadRequest, "error decoding json: %s", err.Error()))
		return
	}

	u := &User{
		Username: params.Username,
		Email:    params.Email,
		password: params.Password,
	}

	req := &CreateUserRequest{
		Interface: httpApiInterface,
		User:      u,
	}

	ExecRequest(w, req)
}

// confirm a user's email address
// func UserConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
// 	u := NewUser(stringParam(ctx, "id"))
// 	u.emailConfirmed = true
// 	logger.Println(u)
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
// func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
// 	u := &User{}
// 	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 		ErrRes(w, NewFmtError(http.StatusBadRequest, err.Error()))
// 		return
// 	}

// 	req := &SaveUserRequest{
// 		// Interface: httpApiInterface,
// 		User:    sessionUser(ctx),
// 		Subject: u,
// 	}

// 	ExecRequest(w, req)
// }

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
