package users

import (
	"database/sql"
	"fmt"
	"github.com/archivers-space/sqlutil"
)

// grab all users
func ReadUsers(db sqlutil.Queryable, limit, offset int) (users []*User, err error) {
	users = make([]*User, 0)
	rows, e := db.Query(fmt.Sprintf("SELECT %s FROM users WHERE deleted=false ORDER BY created DESC LIMIT $1 OFFSET $2", userColumns()), limit, offset)
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

func UsersSearch(db sqlutil.Queryable, query string, limit, offset int) ([]*User, error) {
	q := fmt.Sprintf("%%%s%%", query)
	rows, err := db.Query(qUsersSearch, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}
