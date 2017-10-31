package models

import (
	"fmt"
	"net/http"
	"WorkingPromo/Utils"
)

//Это типа метод который на http запрос отвечает тем что есть в таблице userinfo, так же для теста использовал
// чисто посмотреть что кого, так же через rows.next ходит по строкам, как в джаве вообщем.
func Prom(w http.ResponseWriter) {
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		fmt.Print("ear")
	}
	defer rows.Close()


	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		Utils.CheckError(err)
		fmt.Fprintln(w, "Hi there, I love ", username)
	}
	if err = rows.Err(); err != nil {
		fmt.Print("error")
	}
}
