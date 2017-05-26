package main

import (
	"strings"
)

// UsersRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type UsersRequest struct {
	Interface string
	// the user performing the request
	User *User `required:"true"`
	// users requests embeds pagination info
	Page
}

func (r *UsersRequest) Exec() (interface{}, error) {
	return ReadUsers(appDB, r.Page)
}

type UserRequest struct {
	Interface string
	User      *User
	Subject   *User
}

func (r *UserRequest) Exec() (interface{}, error) {
	if err := r.Subject.Read(appDB); err != nil {
		return nil, err
	}
	return r.Subject, nil
}

type CreateUserRequest struct {
	Interface string
	User      *User
}

func (r *CreateUserRequest) Exec() (interface{}, error) {
	if strings.TrimSpace(r.User.password) == "" {
		return nil, ErrPasswordRequired
	}
	if strings.TrimSpace(r.User.Email) == "" {
		return nil, ErrEmailRequired
	}

	if err := r.User.Save(appDB); err != nil {
		return nil, err
	}
	return r.User, nil
}

type SaveUserRequest struct {
	Interface string
	User      *User
	Subject   *User
}

func (r *SaveUserRequest) Exec() (interface{}, error) {
	if !r.User.isAdmin && r.User.Id != r.Subject.Id {
		return nil, ErrAccessDenied
	}

	log.Info(r.User.Id, r.Subject.Id, r.User.Id == r.Subject.Id)
	if err := r.Subject.Save(appDB); err != nil {
		return nil, err
	}

	return r.Subject, nil
}

type UsersSearchRequest struct {
	Interface string
	User      *User
	Query     string
	Page
}

func (r *UsersSearchRequest) Exec() (interface{}, error) {
	return UsersSearch(appDB, r.Query, r.Page.Size, r.Page.Offset())
}
