package database

import (
	"database/sql"
	"os"
	"path/filepath"
)

func Open(path string) error {
	if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return e
	}
	var e error
	DB, e = sql.Open("sqlite", path)
	if e != nil {
		return e
	}
	DB.SetMaxOpenConns(1)
	Exec("PRAGMA foreign_keys=ON")
	Exec("PRAGMA journal_mode=WAL")
	Exec("PRAGMA busy_timeout=30000")
	return nil
}
