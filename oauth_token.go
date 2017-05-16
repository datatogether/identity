package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/oauth2"
	"time"
)

type UserOauthToken struct {
	User    *User
	Service string
	token   *oauth2.Token
}

func (t *UserOauthToken) ReadUser(db *sql.DB) (*User, error) {
	row, err := db.Query(qUserOauthTokenUser, t.Service, t.token)
	if err != nil {
		return nil, err
	}

	u := &User{}
	if err := u.UnmarshalSQL(row); err != nil {
		return nil, err
	}

	return u, nil
}

func (t *UserOauthToken) UserService() (OauthUserService, error) {
	switch t.Service {
	case OauthServiceGithub:
		return NewGithub(t.token), nil
	default:
		return nil, fmt.Errorf("invalid service name")
	}
}

func (g *UserOauthToken) Read(db *sql.DB) error {
	if g.User == nil || g.Service == "" {
		return ErrNotFound
	}
	return g.UnmarshalSQL(db.QueryRow(qUserOauthTokenByUserAndService, g.User.Id, g.Service))
}

func (t *UserOauthToken) Save(db *sql.DB) error {
	prev := &UserOauthToken{User: t.User, Service: t.Service}
	if err := prev.Read(db); err != nil {
		if err == ErrNotFound {
			if _, err := db.Exec(qUserOauthTokenInsert, t.SQLArgs()...); err != nil {
				logger.Println(err.Error())
				return NewFmtError(500, err.Error())
			}
		} else {
			logger.Println(err.Error())
			return err
		}
	} else {
		if _, err := db.Exec(qUserOauthTokenUpdate, t.SQLArgs()...); err != nil {
			logger.Println(err.Error())
			return err
		}
	}
	return nil
}

func (g *UserOauthToken) Delete(db *sql.DB) error {
	if g.User == nil {
		return ErrNotFound
	}
	_, err := db.Exec(qUserOauthTokenDelete, g.User.Id, g.Service)
	return err
}

// turn a sql.Row result into a reset token pointer
func (g *UserOauthToken) UnmarshalSQL(row sqlScannable) error {
	var (
		userId, service, access, tType, refresh string
		expiry                                  time.Time
	)

	if err := row.Scan(&userId, &service, &access, &tType, &refresh, &expiry); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	tok := &UserOauthToken{
		User:    &User{Id: userId},
		Service: service,
		token: &oauth2.Token{
			AccessToken:  access,
			TokenType:    tType,
			RefreshToken: refresh,
			Expiry:       expiry,
		},
	}

	*g = *tok
	return nil
}

func (t *UserOauthToken) SQLArgs() []interface{} {
	if t.token == nil {
		t.token = &oauth2.Token{}
	}
	return []interface{}{
		t.User.Id,
		t.Service,
		t.token.AccessToken,
		t.token.Type(),
		t.token.RefreshToken,
		t.token.Expiry,
	}
}
