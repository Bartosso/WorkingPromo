package models

import (
	"fmt"
)

func Prom() {
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		fmt.Print("e")
	}
	defer rows.Close()


	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		checkErra(err)
		fmt.Println("hello world " + username)
	}
	if err = rows.Err(); err != nil {
		fmt.Print("ll")
	}
	fmt.Print("ll")
}
func checkErra(err error) {
	if err != nil {
		panic(err)
	}
}