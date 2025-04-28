package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test_tec")

	if err != nil {
		log.Fatal(err)
	}
}
