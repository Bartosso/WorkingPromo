package main

import (
	"fmt"
	"net/http"
	_"github.com/lib/pq"
	"WorkingPromo/Models"
	"WorkingPromo/Utils"

)

//Единственный недохендлер, юзаю для тестов
func handler(w http.ResponseWriter, r *http.Request) {
	//По переходу на данный хендлер (т.е. localhost:8080/) запускается функция - обнови данные в бд по большим папкам
	//типа там всякие игры, кредитки, мобильная связь и т.д.
	//Теперь обновляет еще и игры
	models.UpdateBigFoldersAndGamesFromSections()

}

func main() {
	models.InitDB(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		Utils.DB_USER, Utils.DB_PASSWORD, Utils.DB_NAME))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}



