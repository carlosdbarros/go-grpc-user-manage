package database

import "database/sql"

func initSqliteInMemory() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// 20240909-add-user-table
	userStmt, err := db.Prepare("create table if not exists users (id text, name text, email text, password text)")
	if err != nil {
		return nil, err
	}
	_, err = userStmt.Exec()
	if err != nil {
		return nil, err
	}

	// 20240910-add-permission-table
	permissionStmt, err := db.Prepare("create table if not exists permissions (id text, name text, codename text)")
	if err != nil {
		return nil, err
	}
	_, err = permissionStmt.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}
