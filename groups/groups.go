package groups

import (
	"database/sql"
	"github.com/archivers-space/sqlutil"
)

func ListGroups(db sqlutil.Queryable, limit, offset int) ([]*Group, error) {
	rows, err := db.Query(qGroups, limit, offset)
	if err != nil {
		return nil, err
	}
	return UnmarshalBoundedSqlGroups(rows, limit)
}

func UnmarshalBoundedSqlGroups(rows *sql.Rows, limit int) ([]*Group, error) {
	defer rows.Close()
	gs := make([]*Group, limit)
	i := 0
	for rows.Next() {
		g := &Group{}
		if err := g.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		gs[i] = g
		i++
	}

	return gs[:i], nil
}
