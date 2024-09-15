package configs

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	userStmt, err := db.Prepare("create table if not exists users (id text, name text, email text, password text)")
	if err != nil {
		panic(err)
	}
	_, err = userStmt.Exec()
	if err != nil {
		panic(err)
	}

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
