package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	fmt.Print("lelo")
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

