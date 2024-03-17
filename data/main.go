package data

import (
	"database/sql"
)

func ConnectDatabase() (*sql.DB, error) {
	return sql.Open("sqlite", "./file-system-service-oauth.sqlite")
}

func InitDatabase() (*sql.DB, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}

	//ENUM('classic', 'refresh')
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, IP TEXT NOT NULL, token TEXT, signature TEXT, type TEXT NOT NULL DEFAULT 'classic', active BOOLEAN NOT NULL DEFAULT true, created_at DATETIME NOT NULL)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS signatures (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, signature TEXT, active BOOLEAN DEFAULT true)")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ReadRows[T any](rows *sql.Rows, scan func(t *T) error) (results []T, _err error) {
	_err = nil
	for rows.Next() {
		var line T
		if err := scan(&line); err != nil {
			_err = err
			break
		}
		results = append(results, line)
	}
	return results, _err
}
