package main

import (
	"database/sql"
	"fmt"
)

// grab all users
func ReadUsers(db sqlQueryable, p Page) (users []*User, err error) {
	users = make([]*User, 0)
	rows, e := db.Query(fmt.Sprintf("SELECT %s FROM users WHERE deleted=false ORDER BY created DESC LIMIT $1 OFFSET $2", userColumns()), p.Size, p.Offset())
	if e != nil {
		if e == sql.ErrNoRows {
			return []*User{}, nil
		}
		return nil, New500Error(e.Error())
	}
	defer rows.Close()
	if us, e := scanUsers(rows); e != nil {
		return nil, New500Error(e.Error())
	} else {
		return us, nil
	}
}

// scan a table of users results into a slice of user pointers.
func scanUsers(rows *sql.Rows) ([]*User, error) {
	us := make([]*User, 0)
	for rows.Next() {
		u := &User{}
		err := u.UnmarshalSQL(rows)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}
