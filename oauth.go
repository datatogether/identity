package main

import (
	"database/sql"
	"golang.org/x/oauth2"
)

var (
	githubOAuth *oauth2.Config
)

type OauthToken struct {
	User    *User
	Service string
	Token   string
}

func (g *OauthToken) Read(db *sql.DB) error {
	if g.User == nil || g.Service == "" {
		return ErrNotFound
	}
	return g.UnmarshalSQL(db.QueryRow(qOauthTokenByUserAndService, g.User.Id, g.Service))
}

func (g *OauthToken) Save(db *sql.DB) error {
	prev := &OauthToken{User: g.User, Service: g.Service}
	if err := prev.Read(db); err != nil {
		if err == ErrNotFound {
			if _, err := db.Exec(qOauthTokenInsert, g.SQLArgs()...); err != nil {
				return NewFmtError(500, err.Error())
			}
		} else {
			return err
		}
	} else {
		if _, err := db.Exec(qOauthTokenUpdate, g.SQLArgs()...); err != nil {
			return err
		}
	}
	return nil
}

func (g *OauthToken) Delete(db *sql.DB) error {
	if g.User == nil {
		return ErrNotFound
	}
	_, err := db.Exec(qOauthTokenDelete, g.User.Id, g.Service)
	return err
}

// turn a sql.Row result into a reset token pointer
func (g *OauthToken) UnmarshalSQL(row sqlScannable) error {
	var (
		userId, service, token string
	)

	if err := row.Scan(&userId, &service, &token); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	tok := &OauthToken{
		User:    &User{Id: userId},
		Service: service,
		Token:   token,
	}

	*g = *tok
	return nil
}

func (g *OauthToken) SQLArgs() []interface{} {
	return []interface{}{
		g.User.Id,
		g.Service,
		g.Token,
	}
}
