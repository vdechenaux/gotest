package database

import (
	"database/sql"
	"sync"
)

var DB *sql.DB
var Mutex sync.Mutex

func Open(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)

	return err
}
