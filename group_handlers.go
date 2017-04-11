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
	req := &GroupRequest{
		Group: &Group{
			Id: r.FormValue("id"),
		},
	}
	ExecRequest(w, envelope, req)
}

// list users or get a single user if supplied with a "username" formValue
func ListGroupsHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	envelope := r.FormValue("envelope") != "false"

	req = &GroupsRequest{
		Interface: httpApiInterface,
		User:      sessionUser(r),
		Page:      PageFromRequest(r),
	}

	ExecRequest(w, envelope, req)
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
	req := &SaveGroupRequest{
		Interface: httpApiInterface,
		User:      sess,
		Group:     g,
	}

	ExecRequest(w, true, req)
}

// delete a user
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	req := &DeleteGroupRequest{
		User: sessionUser(r),
		Group: &Group{
			Id: r.FormValue("id"),
		},
	}

	ExecRequest(w, true, req)
}
