package main

import (
	"encoding/json"
	"github.com/archivers-space/identity/groups"
	"net/http"
)

var GroupsRequests = groups.GroupRequests{Store: appDB}

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
	p := &groups.GroupsGetParams{
		Group: &groups.Group{
			Id: r.FormValue("id"),
		},
	}

	res := &groups.Group{}
	if err := GroupsRequests.Get(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, envelope, res)
}

// list users or get a single user if supplied with a "username" formValue
func ListGroupsHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"

	page := PageFromRequest(r)
	p := &groups.GroupsListParams{
		User:   sessionUser(r),
		Limit:  page.Size,
		Offset: page.Offset(),
	}

	res := []*groups.Group{}
	if err := GroupsRequests.List(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, envelope, res)
}

// SaveGroupHandler updates a user
func SaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessionUser(r)
	g := &groups.Group{}
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
	p := &groups.GroupsSaveParams{
		User:  sess,
		Group: g,
	}
	res := &groups.Group{}
	if err := GroupsRequests.Save(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

// delete a user
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	var res bool
	p := &groups.GroupsDeleteParams{
		User: sessionUser(r),
		Group: &groups.Group{
			Id: r.FormValue("id"),
		},
	}

	if err := GroupsRequests.Delete(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}
