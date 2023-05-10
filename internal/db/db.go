package db

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

func Migrate(connStr string) {
	rootDir, _ := os.Getwd()
	tables, _ := ioutil.ReadFile(rootDir + "/migrations/tables.sql")
	constraints, _ := ioutil.ReadFile(rootDir + "/migrations/constraints.sql")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(string(tables))
	if err != nil {
		panic(err)
	}
	_, _ = db.Exec(string(constraints))
}
