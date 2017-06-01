package main

import (
	"strings"
)

// Users holds all types of requests for users
// it's based on an int b/c it's stateless and Go lets us
// do this sort of thing
type Users int

// UsersRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type UsersListParams struct {
	// the user performing the request
	User *User `required:"true"`
	// users requests embeds pagination info
	Page
}

func (Users) List(p *UsersListParams, res *[]*User) error {
	users, err := ReadUsers(appDB, p.Page)
	if err != nil {
		return err
	}

	*res = users
	return nil
}

type UsersGetParams struct {
	User    *User
	Subject *User
}

func (Users) Get(p *UsersGetParams, res *User) error {
	if err := p.Subject.Read(appDB); err != nil {
		return err
	}

	*res = *p.Subject
	return nil
}

type UsersCreateParams struct {
	User *User
}

func (Users) Create(p *UsersCreateParams, res *User) error {
	if strings.TrimSpace(p.User.password) == "" {
		return ErrPasswordRequired
	}
	if strings.TrimSpace(p.User.Email) == "" {
		return ErrEmailRequired
	}

	if err := p.User.Save(appDB); err != nil {
		return err
	}

	*res = *p.User
	return nil
}

type UsersSaveParams struct {
	User    *User
	Subject *User
}

func (Users) Save(p *UsersSaveParams, res *User) error {
	if !p.User.isAdmin && p.User.Id != p.Subject.Id {
		return ErrAccessDenied
	}

	if err := p.Subject.Save(appDB); err != nil {
		return err
	}

	*res = *p.Subject
	return nil
}

type UsersSearchParams struct {
	User  *User
	Query string
	Page
}

func (Users) Search(p *UsersSearchParams, res *[]*User) error {
	users, err := UsersSearch(appDB, p.Query, p.Page.Size, p.Page.Offset())
	if err != nil {
		return err
	}

	*res = users
	return nil
}
