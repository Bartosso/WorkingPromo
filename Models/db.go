package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
)
//Глобальная переменная для доступа к бд, тот самый пул
var db *sql.DB
//Метод что создает базу данных и проверяет получил ли в ответ ошибку.
//В конце пишет что датабеис готово, на самом деле криво пздц
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	fmt.Println("Database ready!")
}

