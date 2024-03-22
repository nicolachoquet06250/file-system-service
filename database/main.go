package database

import (
	"database/sql"
	"filesystem_service/flags"
	"fmt"
	"os"
	"strings"
)

var databaseFileName = "file-system-service-oauth.sqlite"

func checkFatalPermissionDenied(err error) {
	if strings.Contains(err.Error(), "readonly database") {
		fmt.Printf("La base de données est en lecture seule car vous n'avez pas les droit d'écriture.\n")
		fmt.Printf("Lancer le service en sudo pour avoir les droits d'écriture.\n\n")
		fmt.Printf("sudo %v\n", strings.Join(os.Args, " "))
		os.Exit(2)
	}
}

func Connect() (*sql.DB, error) {
	path := "."
	if flags.IsProd() {
		path = prodPath
	}

	return sql.Open("sqlite", path+"/"+databaseFileName)
}

func Init() (*sql.DB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS role_actions (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		role_action_name VARCHAR(255)
	);`); err != nil {
		return nil, err
	}
	if _, err = db.Exec(`INSERT OR IGNORE INTO role_actions (id, role_action_name) VALUES
                                                              (1, 'read_dir'),
                                                              (2, 'create_dir'),
                                                              (3, 'delete_dir'),
                                                              (4, 'rename_dir'),
                                                              (5, 'read_file'),
                                                              (6, 'create_file'),
                                                              (7, 'delete_file'),
                                                              (8, 'rename_file'),
                                                              (9, 'update_file_content');`); err != nil {
		checkFatalPermissionDenied(err)
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		role_name VARCHAR(255) NOT NULL,
		active BOOLEAN NOT NULL DEFAULT FALSE
	);`); err != nil {
		return nil, err
	}
	if _, err = db.Exec(`INSERT OR IGNORE INTO roles (id, role_name, active) VALUES
                                                        (1, 'readwrite', TRUE),
                                                        (2, 'readonly', TRUE);`); err != nil {
		checkFatalPermissionDenied(err)
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS roles_link_role_actions (
		role INTEGER NOT NULL,
		role_action INTEGER NOT NULL,
		PRIMARY KEY (role, role_action),
		FOREIGN KEY (role) REFERENCES roles(id),
		FOREIGN KEY (role_action) REFERENCES role_actions(id)
	);`); err != nil {
		return nil, err
	}
	if _, err = db.Exec(`INSERT OR IGNORE INTO roles_link_role_actions (role, role_action) VALUES
                                                        (1, 1),
                                                        (1, 2),
                                                        (1, 3),
                                                        (1, 4),
                                                        (1, 5),
                                                        (1, 6),
                                                        (1, 7),
                                                        (1, 8),
                                                        (1, 9),
                                                        (2, 1),
                                                        (2, 5);`); err != nil {
		checkFatalPermissionDenied(err)
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS credentials (
		client_id VARCHAR(255) PRIMARY KEY,
		client_secret VARCHAR(255) NOT NULL,
		role INTEGER NOT NULL,
		creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_date TIMESTAMP DEFAULT NULL,
-- 		1 an (en timestamp)
		expires_in INTEGER DEFAULT 31536000,
		active BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (role) REFERENCES roles(id)
	)`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	ip VARCHAR(255) NOT NULL,
		client_id VARCHAR(255),
		access_token VARCHAR(255) NOT NULL,
		refresh_token VARCHAR(255) NOT NULL,
-- 		1h (en timestamp)
		expires_in INT DEFAULT NULL DEFAULT 3600,
		creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		active BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (client_id) REFERENCES credentials(client_id)
	)`); err != nil {
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
