package main

import (
	"fmt"
	"net/http"
	_"github.com/lib/pq"
	"WorkingPromo/Models"
	"WorkingPromo/Utils"
	"github.com/gorilla/mux"


)

//Единственный недохендлер, юзаю для тестов
func handler(w http.ResponseWriter, r *http.Request) {
	//По переходу на данный хендлер (т.е. localhost:8080/) запускается функция - обнови данные в бд по большим папкам
	//типа там всякие игры, кредитки, мобильная связь и т.д.
	//Теперь обновляет еще и игры с продовцами
	models.UpdateBigFoldersAndGamesFromSections()
	models.UpdateGamesOffersAndSellersStats()

}

func getSectionsHandler(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-type", "application/json")
	chosenID := r.URL.Query().Get("id")
	if chosenID == "0" {
		w.Write(*models.GetAvailableGamesSectionsAndFoldersViaJson())
	} else {
		result := models.GetSelectedByParentIDSections(chosenID)
		if result == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No result"))
		} else {
			w.Write(*result)
		}

	}

}


func getOffersHandler(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-type", "application/json")
	chosenID := r.URL.Query().Get("id")
		result := models.GetOffersViaJson(chosenID)
		if result == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No result"))
		} else {
			w.Write(*result)
		}



}

func main() {
	models.InitDB(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		Utils.DB_USER, Utils.DB_PASSWORD, Utils.DB_NAME))

	rtr := mux.NewRouter()
	rtr.HandleFunc("/getoffers/", getOffersHandler).
		Methods("GET").Queries("id", "{id:[0-9]+}")
	rtr.HandleFunc("/getsections/", getSectionsHandler).
		Methods("GET").Queries("id", "{id:[0-9]+}")
	http.Handle("/", rtr)


	//http.HandleFunc("/update", handler)
	//http.HandleFunc("/getsections/{id:[0-9]}", getSectionsHandler)
	http.ListenAndServe(":8080", nil)
}



