package main

import (
	"fmt"
	"net/http"
	_"github.com/lib/pq"


)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "test"
	DB_NAME     = "test"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("lel")
	.Prom()
}

func main() {

	.InitDB(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME))

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

