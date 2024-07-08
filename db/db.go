package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewSQLDatabase(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT version()")

	if err != nil {
		db.Close()
		return nil, err
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("version=%s\n", version)

	return db, nil
}
