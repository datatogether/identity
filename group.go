package main

import (
	"database/sql"
	"time"
)

// a user reset token
type Group struct {
	Id          string `json:"id" sql:"id"`
	Created     int64  `json:"created" sql:"created"`
	Updated     int64  `json:"updated" sql:"updated"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	PosterUrl   string `json:"posterUrl"`
	ProfileUrl  string `json:"profileUrl"`
	Creator     *User  `json:"creator"`
}

// validate a group
func (r *Group) validate(db *sql.DB) error {
	return nil
}

func (g *Group) InviteUser(db *sql.DB, u *User) error {
	_, err := db.Exec(qGroupInviteUser, g.Id, u.Id)
	return err
}

// read a group
func (g *Group) Read(db *sql.DB) error {
	if g.Id == "" {
		return ErrNotFound
	}
	return g.UnmarshalSQL(db.QueryRow(qGroupById, g.Id))
}

func (g *Group) Save(db *sql.DB) error {
	prev := &Group{Id: g.Id}
	if err := prev.Read(db); err != nil {
		if err == ErrNotFound {
			g.Id = NewUuid()
			g.Created = time.Now().Unix()
			g.Updated = g.Created
			if _, err := db.Exec(qGroupInsert, g.SQLArgs()...); err != nil {
				return NewFmtError(500, err.Error())
			}
		} else {
			return err
		}
	} else {
		g.Updated = time.Now().Unix()
		if _, err := db.Exec(qGroupUpdate, g.SQLArgs()...); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) Delete(db *sql.DB) error {
	_, err := db.Exec(qGroupDelete, g.Id)
	return err
}

// turn a sql.Row result into a reset token pointer
func (g *Group) UnmarshalSQL(row sqlScannable) error {
	var (
		id, title, description, color, posterUrl, profileUrl, creatorId string
		created, updated                                                int64
	)

	if err := row.Scan(&id, &created, &updated, &title, &description, &color, &profileUrl, &posterUrl, &creatorId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	group := &Group{
		Id:          id,
		Created:     created,
		Updated:     updated,
		Title:       title,
		Description: description,
		Color:       color,
		ProfileUrl:  profileUrl,
		PosterUrl:   posterUrl,
		Creator:     &User{Id: creatorId},
	}

	*g = *group
	return nil
}

func (g *Group) SQLArgs() []interface{} {
	return []interface{}{
		g.Id,
		g.Created,
		g.Updated,
		g.Title,
		g.Description,
		g.Color,
		g.ProfileUrl,
		g.PosterUrl,
		g.Creator.Id,
	}
}
