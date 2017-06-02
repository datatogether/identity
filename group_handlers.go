package main

import (
	"encoding/json"
	"net/http"
)

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		ListGroupsHandler(w, r)
	case "POST", "PUT":
		SaveGroupHandler(w, r)
	case "DELETE":
		DeleteGroupHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		ReadGroupHandler(w, r)
	case "POST", "PUT":
		SaveGroupHandler(w, r)
	case "DELETE":
		DeleteGroupHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func ReadGroupHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"
	p := &GroupsGetParams{
		Group: &Group{
			Id: r.FormValue("id"),
		},
	}

	res := &Group{}
	if err := new(Groups).Get(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, envelope, res)
}

// list users or get a single user if supplied with a "username" formValue
func ListGroupsHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"

	p := &GroupsListParams{
		User: sessionUser(r),
		Page: PageFromRequest(r),
	}

	res := []*Group{}
	if err := new(Groups).List(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, envelope, res)
}

// SaveGroupHandler updates a user
func SaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessionUser(r)
	g := &Group{}
	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, err.Error()))
			return
		}
	} else {
		g.Title = r.FormValue("title")
		g.Description = r.FormValue("description")
		g.Color = r.FormValue("color")
		g.ProfileUrl = r.FormValue("profileUrl")
		g.PosterUrl = r.FormValue("posterUrl")
	}

	g.Creator = sess
	p := &GroupsSaveParams{
		User:  sess,
		Group: g,
	}
	res := &Group{}
	if err := new(Groups).Save(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

// delete a user
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	var res bool
	p := &GroupsDeleteParams{
		User: sessionUser(r),
		Group: &Group{
			Id: r.FormValue("id"),
		},
	}

	if err := new(Groups).Delete(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}
