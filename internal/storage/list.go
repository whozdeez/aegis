package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type EntrySummary struct {
	Service  string
	Username string
}

func ListEntries(dbPath string) ([]EntrySummary, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT service, username FROM entries ORDER BY service`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []EntrySummary
	for rows.Next() {
		var e EntrySummary
		if err := rows.Scan(&e.Service, &e.Username); err != nil {
			return nil, err
		}
		results = append(results, e)
	}

	return results, nil
}
