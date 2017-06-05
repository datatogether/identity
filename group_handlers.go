package main

import (
	"encoding/json"
	"github.com/archivers-space/identity/group"
	"net/http"
)

var GroupsRequests = group.GroupRequests{Store: appDB}

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
	p := &group.GroupsGetParams{
		Group: &group.Group{
			Id: r.FormValue("id"),
		},
	}

	res := &group.Group{}
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
	p := &group.GroupsListParams{
		User:   sessionUser(r),
		Limit:  page.Size,
		Offset: page.Offset(),
	}

	res := []*group.Group{}
	if err := GroupsRequests.List(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, envelope, res)
}

// SaveGroupHandler updates a user
func SaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessionUser(r)
	g := &group.Group{}
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
	p := &group.GroupsSaveParams{
		User:  sess,
		Group: g,
	}
	res := &group.Group{}
	if err := GroupsRequests.Save(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

// delete a user
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	var res bool
	p := &group.GroupsDeleteParams{
		User: sessionUser(r),
		Group: &group.Group{
			Id: r.FormValue("id"),
		},
	}

	if err := GroupsRequests.Delete(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}
